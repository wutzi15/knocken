package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/profile"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wutzi15/knocken/config"
	containscheck "github.com/wutzi15/knocken/containsCheck"
	diffcheck "github.com/wutzi15/knocken/diffCheck"
	"github.com/wutzi15/knocken/parsers"
	types "github.com/wutzi15/knocken/types"
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

	wpPosts = prometheus.NewGaugeVec(
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
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	// defer profile.Start(profile.CPUProfile).Stop()
	config := config.GetConfig()

	verbose = config.Verbose
	waitTime := config.WaitTime

	URLs, err := parsers.ParseTargets(config.Targets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("URLS to check: %+v\n", URLs)

	ignore, err := parsers.ParseTargets(config.Ignore)

	if err == nil {
		fmt.Printf("URLS to ignore: %+v\n", ignore)
		URLs = parsers.RemoveIgnoredTargets(URLs, ignore)
	}

	contains, err := parsers.ParseContainsTargets(config.ContainsTargets)

	if err == nil {
		fmt.Printf("URLS to check for contains: %+v\n", contains)
	}

	fmt.Printf("URLS to check for contains: %+v\n", contains)

	prometheus.MustRegister(statSame)
	prometheus.MustRegister(statContains)

	cfg := types.MetricsConfig{
		URLs:     URLs,
		SaveDiff: config.SaveDiff,
		FastDiff: config.FastDiff,
		WaitTime: waitTime,
		StatSame: statSame,
		Verbose:  verbose,
		Wg:       nil,
	}
	if config.RunDiff {
		go func() {
			diffcheck.RecordMetrics(cfg)
		}()
	}

	containcfg := types.ContainsConfig{
		WaitTime:     waitTime,
		Verbose:      verbose,
		StatContains: statContains,
		Wg:           nil,
	}
	if config.RunContain {
		containscheck.RunContain(contains, containcfg)
	}
	// recordMetrics(URLs, *saveDiff, waitTime)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
