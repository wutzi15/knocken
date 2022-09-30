package diffcheck

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"strings"
	"time"

	myTypes "github.com/wutzi15/knocken/types"
	"github.com/wutzi15/knocken/utils"
	"github.com/wutzi15/levenshtein"
)

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
			for _, target := range config.URLs.Targets {
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
				htmlNew, err := utils.GetHTML(target)
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
				same := 0.0
				// diffs := levenshtein.Distance(htmlNewStr, htmlOldStr, false)
				if !config.FastDiff {
					levenshteinDiff := float64(levenshtein.Distance(htmlNewStr, htmlOldStr))
					fmt.Println("levenshteinDiff: ", levenshteinDiff)
					len1 := float64(len(htmlNewStr))
					len2 := float64(len(htmlOldStr))
					weightedLen := (len1 + len2) / 2.0
					same = math.Abs(1 - (levenshteinDiff / weightedLen))
				} else {
					same = levenshtein.DistanceJW(htmlNewStr, htmlOldStr)
				}
				if config.Verbose {
					fmt.Printf("Target: %s\n\tSame: %f\n", target, same)
				}
				if config.SaveDiff {
					WriteFile("same_"+hostname, []byte(fmt.Sprintf("%e", same)))
				}
				fmt.Println(same)
				config.StatSame.WithLabelValues(target).Set(same)
			}
			if config.Wg != nil {
				config.Wg.Done()
			}
			time.Sleep(config.WaitTime)
		}
	}()
}
