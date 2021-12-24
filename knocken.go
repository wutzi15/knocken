package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func recordMetrics(URLs URL) {
	go func() {
		for {
			// Read file to byte array
			// data, err := ioutil.ReadFile("targets.yml")
			// if err != nil {
			// 	panic(err)
			// }

			for _, url := range URLs {
				for _, target := range url.Targets {
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
					same := 1 - (levenshteinDiff / weightedLen)
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

	data, err := ioutil.ReadFile("targets.yml")
	if err != nil {
		panic(err)
	}

	var URLs URL
	err = yaml.Unmarshal(data, &URLs)
	if err != nil {
		panic(err)
	}

	var ignore URL
	ignoreData, err := ioutil.ReadFile("ignore.yml")
	if err == nil {
		err = yaml.Unmarshal(ignoreData, &ignore)
		if err != nil {
			panic(err)
		}
		// Really ugly nested loops, but we need to stick to the YAML format giben by prometheus
		for _, ign := range ignore {
			for _, ignUrl := range ign.Targets {
				for _, url := range URLs {
					for l, urlTarget := range url.Targets {
						if urlTarget == ignUrl {
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
	recordMetrics(URLs)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
