package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wutzi15/knocken/config"
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
			Help:      "Percentage of same HTML code on a website in the last 5 min.",
		},
		[]string{
			"target",
		},
	)
	verbose = false
)

func main() {
	fmt.Println("Starting...")

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

	prometheus.MustRegister(statSame)
	cfg := types.MetricsConfig{
		URLs:     URLs,
		SaveDiff: config.SaveDiff,
		WaitTime: waitTime,
		StatSame: statSame,
		Verbose:  verbose,
		Wg:       nil,
	}
	diffcheck.RecordMetrics(cfg)
	// recordMetrics(URLs, *saveDiff, waitTime)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
