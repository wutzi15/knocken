package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wutzi15/knocken/config"
	diffcheck "github.com/wutzi15/knocken/diffCheck"
	types "github.com/wutzi15/knocken/types"
	"gopkg.in/yaml.v2"
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

	data, err := ioutil.ReadFile(config.Targets)
	if err != nil {
		panic(err)
	}

	var URLs types.URL
	err = yaml.Unmarshal(data, &URLs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("URLS to check: %+v\n", URLs)

	var ignore types.URL
	ignoreData, err := ioutil.ReadFile(config.Ignore)
	if err == nil {
		err = yaml.Unmarshal(ignoreData, &ignore)
		if err != nil {
			panic(err)
		}

		fmt.Printf("URLS to ignore: %+v\n", ignore)

		// Really ugly nested loops, but we need to stick to the YAML format giben by prometheus
		for _, ign := range ignore {
			for _, ignUrl := range ign.Targets {
				for _, url := range URLs {
					for l, urlTarget := range url.Targets {
						if urlTarget == ignUrl {
							if verbose {
								fmt.Printf("Ignoring %s\n", urlTarget)
							}
							copy(url.Targets[l:], url.Targets[l+1:])
							url.Targets[len(url.Targets)-1] = ""
							url.Targets = url.Targets[:len(url.Targets)-1]
						}
					}
				}
			}
		}
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
