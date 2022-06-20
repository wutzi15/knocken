package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v2"
)

type URL []struct {
	Targets []string `yaml:"targets"`
}

var (
	statSame = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "imflow",
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

func getHTML(url string) ([]byte, error) {
	resp, err := http.Get("https://" + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func writeFile(fileName string, data []byte) error {
	return ioutil.WriteFile("./html/"+fileName, data, 0644)
}

func getContentOfFileIfExists(fileName string) ([]byte, error) {
	if _, err := ioutil.ReadFile("./html/" + fileName); err != nil {
		return nil, err
	}
	return ioutil.ReadFile("./html/" + fileName)
}

func recordMetrics(URLs URL, saveDiff bool) {
	go func() {
		for {
			// Read file to byte array
			// data, err := ioutil.ReadFile("targets.yml")
			// if err != nil {
			// 	panic(err)
			// }

			for _, url := range URLs {
				for _, target := range url.Targets {
					if strings.TrimSpace(target) == "" {
						continue
					}
					fmt.Println("Checking: " + target)
					dmp := diffmatchpatch.New()
					htmlNew, err := getHTML(target)
					if err != nil {
						panic(err)
					}
					htmlOld, err := getContentOfFileIfExists(target)

					writeFile(target, htmlNew)
					if err != nil {
						continue
					}
					htmlNewStr := string(htmlNew)
					htmlOldStr := string(htmlOld)
					diffs := dmp.DiffMain(htmlNewStr, htmlOldStr, false)
					levenshteinDiff := float64(dmp.DiffLevenshtein(diffs))
					len1 := float64(len(htmlNewStr))
					len2 := float64(len(htmlOldStr))
					weightedLen := (len1 + len2) / 2.0
					same := math.Abs(1 - (levenshteinDiff / weightedLen))
					if verbose {
						fmt.Printf("\nLevenshtein: %f\nWeightedLen: %f\nSame: %f\n", levenshteinDiff, weightedLen, same)
					}
					if saveDiff {
						writeFile(fmt.Sprint(same)+"_"+target, htmlNew)
					}
					fmt.Println(same)
					statSame.WithLabelValues(target).Set(same)
				}
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func main() {
	// url := "https://imflow.me"
	fmt.Println("Starting...")

	saveDiff := flag.Bool("saveDiffs", false, "Keep diffs in ./html/ with diff percentage")
	v := flag.Bool("verbose", false, "Verbose output")

	flag.Parse()

	verbose = *v

	data, err := ioutil.ReadFile("targets.yml")
	if err != nil {
		panic(err)
	}

	var URLs URL
	err = yaml.Unmarshal(data, &URLs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("URLS to check: %+v\n", URLs)

	var ignore URL
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
	recordMetrics(URLs, *saveDiff)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
