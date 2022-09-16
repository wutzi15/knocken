package diffcheck

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/wutzi15/knocken/types"
	"github.com/wutzi15/levenshtein"
)

func GetHTML(url string) ([]byte, error) {
	var domain string
	if strings.HasPrefix(url, "http") {
		domain = url
	} else {
		domain = "https://" + url
	}
	resp, err := http.Get(domain)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func WriteFile(fileName string, data []byte) error {
	//check if folder exists and create it if not
	if _, err := ioutil.ReadDir("./html"); err != nil {
		if err := os.Mkdir("./html", 0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile("./html/"+fileName, data, 0644)
}

func getContentOfFileIfExists(fileName string) ([]byte, error) {
	if _, err := ioutil.ReadFile("./html/" + fileName); err != nil {
		return nil, err
	}
	return ioutil.ReadFile("./html/" + fileName)
}

func RecordMetrics(URLs types.URL, saveDiff bool, waitTime time.Duration, statSame *prometheus.GaugeVec, verbose bool) {
	go func() {
		for {
			for _, url := range URLs {
				for _, target := range url.Targets {
					if strings.TrimSpace(target) == "" {
						continue
					}
					fmt.Println("Checking: " + target)
					// dmp := diffmatchpatch.New()
					htmlNew, err := GetHTML(target)
					if err != nil {
						panic(err)
					}
					htmlOld, err := getContentOfFileIfExists(target)

					WriteFile(target, htmlNew)
					if err != nil {
						continue
					}
					htmlNewStr := string(htmlNew)
					htmlOldStr := string(htmlOld)
					// diffs := levenshtein.Distance(htmlNewStr, htmlOldStr, false)
					levenshteinDiff := float64(levenshtein.Distance(htmlNewStr, htmlOldStr))
					len1 := float64(len(htmlNewStr))
					len2 := float64(len(htmlOldStr))
					weightedLen := (len1 + len2) / 2.0
					same := math.Abs(1 - (levenshteinDiff / weightedLen))
					if verbose {
						fmt.Printf("\nLevenshtein: %f\nWeightedLen: %f\nSame: %f\n", levenshteinDiff, weightedLen, same)
					}
					if saveDiff {
						WriteFile(fmt.Sprint(same)+"_"+target, htmlNew)
					}
					fmt.Println(same)
					statSame.WithLabelValues(target).Set(same)
				}
			}
			time.Sleep(waitTime)
		}
	}()
}
