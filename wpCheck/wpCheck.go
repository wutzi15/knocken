package wpcheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/wutzi15/knocken/types"
	"github.com/wutzi15/knocken/utils"
)

func WriteJSONFile(fileName string, data types.WPPosts) error {
	//check if folder exists and create it if not
	if _, err := ioutil.ReadDir("./json"); err != nil {
		// fmt.Println("Creating folder ./html")
		if err := os.Mkdir("./json", 0755); err != nil {
			// fmt.Println("Error creating folder ./html", err)
			return err
		}
	}
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("./json/"+fileName, d, 0644)
}

func GetJSONContentOfFileIfExists(fileName string) (types.WPPosts, error) {
	data, err := ioutil.ReadFile("./json/" + fileName)
	if err != nil {
		return nil, err
	}
	var ret types.WPPosts
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func WPPostsFunc(targets types.URL, config types.WPPostsConfig) {
	// var ret []bool
	for _, target := range targets.Targets {
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

		dataNew, err := utils.GetJSON(hostname + "/wp-json/wp/v2/posts?_fields=author,id,date,title")
		if err != nil {
			fmt.Printf("Error getting json from %s: %s)\n", target, err)
		}
		var Posts types.WPPosts
		err = json.Unmarshal(dataNew, &Posts)

		Old, err := GetJSONContentOfFileIfExists(hostname)
		if err != nil {
			fmt.Printf("Error getting json from %s: %s\n)", target, err)
		}
		WriteJSONFile(hostname, Posts)

		fmt.Printf("Length of old: %d, new: %d\n", len(Old), len(Posts))
		timeInHours := config.WaitTime.Hours()
		fmt.Printf("Waittime: %f\n", timeInHours)
		newPerHour := float64(len(Posts)-len(Old)) / timeInHours
		fmt.Printf("New posts per hour: %f\n", newPerHour)
		config.StatWpPosts.WithLabelValues(target).Set(newPerHour)
	}
}

func RunWPChecksPosts(targets types.URL, config types.WPPostsConfig) {
	go func() {
		for {
			WPPostsFunc(targets, config)
			if config.Wg != nil {
				config.Wg.Done()
			}
			time.Sleep(config.WaitTime)
		}
	}()
}
