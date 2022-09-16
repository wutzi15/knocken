package diffcheck

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	myTypes "github.com/wutzi15/knocken/types"
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
		// fmt.Println("Creating folder ./html")
		if err := os.Mkdir("./html", 0755); err != nil {
			// fmt.Println("Error creating folder ./html", err)
			return err
		}
	}
	return ioutil.WriteFile("./html/"+fileName, data, 0644)
}

func GetContentOfFileIfExists(fileName string) ([]byte, error) {
	if _, err := ioutil.ReadFile("./html/" + fileName); err != nil {
		return nil, err
	}
	return ioutil.ReadFile("./html/" + fileName)
}

func RecordMetrics(config myTypes.MetricsConfig) {
	go func() {
		for {
			for _, urls := range config.URLs {
				for _, target := range urls.Targets {
					if strings.TrimSpace(target) == "" {
						continue
					}

					var toUrlToParse = target
					if !strings.HasPrefix(target, "http") {
						toUrlToParse = "https://" + target
					}
					myurl, err := url.ParseRequestURI(toUrlToParse)
					if err != nil {
						fmt.Println(err)
						continue
					}
					hostname := myurl.Hostname()
					// dmp := diffmatchpatch.New()
					htmlNew, err := GetHTML(target)
					if err != nil {
						panic(err)
					}
					htmlOld, err := GetContentOfFileIfExists(hostname)

					if err != nil {
						fmt.Println("Error: " + err.Error())
					}
					WriteFile(hostname, htmlNew)
					htmlNewStr := string(htmlNew)
					htmlOldStr := string(htmlOld)
					// diffs := levenshtein.Distance(htmlNewStr, htmlOldStr, false)
					levenshteinDiff := float64(levenshtein.Distance(htmlNewStr, htmlOldStr))
					len1 := float64(len(htmlNewStr))
					len2 := float64(len(htmlOldStr))
					weightedLen := (len1 + len2) / 2.0
					same := math.Abs(1 - (levenshteinDiff / weightedLen))
					if config.Verbose {
						fmt.Printf("\nLevenshtein: %f\nWeightedLen: %f\nSame: %f\n", levenshteinDiff, weightedLen, same)
					}
					if config.SaveDiff {
						WriteFile("same_"+hostname, []byte(fmt.Sprintf("%e", same)))
					}
					fmt.Println(same)
					config.StatSame.WithLabelValues(target).Set(same)
				}
			}
			if config.Wg != nil {
				config.Wg.Done()
			}
			time.Sleep(config.WaitTime)
		}
	}()
}
