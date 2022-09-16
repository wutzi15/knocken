package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	saveDiff := flag.Bool("saveDiffs", false, "Keep diffs in ./html/ with diff percentage")
	v := flag.Bool("verbose", false, "Verbose output")
	waitTimeStr := flag.String("waitTime", "5m", "Wait time")

	flag.Parse()

	verbose = *v
	waitTime, err := time.ParseDuration(*waitTimeStr)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("targets.yml")
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
	ignoreData, err := ioutil.ReadFile("ignore.yml")
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
	diffcheck.RecordMetrics(URLs, *saveDiff, waitTime, statSame, verbose)
	// recordMetrics(URLs, *saveDiff, waitTime)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
