package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wutzi15/knocken/config"
	containscheck "github.com/wutzi15/knocken/containsCheck"
	diffcheck "github.com/wutzi15/knocken/diffCheck"
	"github.com/wutzi15/knocken/parsers"
	types "github.com/wutzi15/knocken/types"
	wpcheck "github.com/wutzi15/knocken/wpCheck"
)

var (
	statSame = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "knocken",
			Subsystem: "knocken",
			Name:      "same",
			Help:      "Percentage of same HTML code on a website since the last checks.",
		},
		[]string{
			"target",
		},
	)

	statContains = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "knocken",
			Subsystem: "knocken",
			Name:      "contains",
			Help:      "1 if the website contains the string 0 otherwise",
		},
		[]string{
			"target",
		},
	)

	statwpPosts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "knocken",
			Subsystem: "knocken",
			Name:      "wpPosts",
			Help:      "number of new posts per 1 hour on a wordpress website. ",
		},
		[]string{
			"target",
		},
	)
	verbose = false
)

func main() {
	fmt.Println("Starting...")

	// profile cpu and memory
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	// defer profile.Start(profile.CPUProfile).Stop()
	config := config.GetConfig()

	verbose = config.Verbose
	waitTime := config.WaitTime

	prometheus.MustRegister(statSame)
	prometheus.MustRegister(statContains)
	prometheus.MustRegister(statwpPosts)

	if config.RunDiff {
		URLs, err := parsers.ParseTargets(config.Targets)
		if err != nil {
			fmt.Printf("Error parsing targets: %s", err)
		}

		fmt.Printf("URLS to check: %+v\n", URLs)

		ignore, err := parsers.ParseTargets(config.Ignore)

		if err == nil {
			fmt.Printf("URLS to ignore: %+v\n", ignore)
			URLs = parsers.RemoveIgnoredTargets(URLs, ignore)
		}
		cfg := types.MetricsConfig{
			URLs:     URLs,
			SaveDiff: config.SaveDiff,
			FastDiff: config.FastDiff,
			WaitTime: waitTime,
			StatSame: statSame,
			Verbose:  verbose,
			Wg:       nil,
		}
		go func() {
			diffcheck.RecordMetrics(cfg)
		}()
	}

	if config.RunContain {
		containcfg := types.ContainsConfig{
			WaitTime:     waitTime,
			Verbose:      verbose,
			StatContains: statContains,
			Wg:           nil,
		}
		contains, err := parsers.ParseContainsTargets(config.ContainsTargets)

		if err != nil {
			fmt.Printf("Error parsing contains targets: %s", err)
		} else {
			fmt.Printf("URLS to check for contains: %+v\n", contains)
			go func() {
				containscheck.RunContain(contains, containcfg)
			}()
		}

	}

	if config.RunWPPosts {
		wpcfg := types.WPPostsConfig{
			WaitTime:    waitTime,
			Verbose:     verbose,
			StatWpPosts: statwpPosts,
			Wg:          nil,
		}

		wp, err := parsers.ParseTargets(config.WPTargets)

		if err != nil {
			fmt.Printf("Error parsing targets: %s", err)
		} else {
			go func() {
				wpcheck.RunWPChecksPosts(wp, wpcfg)
			}()
		}
	}
	// recordMetrics(URLs, *saveDiff, waitTime)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
