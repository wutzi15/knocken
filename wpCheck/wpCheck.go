package wpcheck

import (
	"bytes"
	"encoding/binary"
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

func writeFile(filename string, data float64) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, data)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return err
	}
	//check if folder exists and create it if not
	if _, err := ioutil.ReadDir("./json"); err != nil {
		// fmt.Println("Creating folder ./json")
		if err := os.Mkdir("./json", 0755); err != nil {
			// fmt.Println("Error creating folder ./json", err)
			return err
		}
	}
	return ioutil.WriteFile("./json/"+filename, buf.Bytes(), 0644)
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

		var dom = ""
		if config.Testing {
			dom = target
		} else {
			dom = target + "/wp-json/wp/v2/posts?_fields=author,id,date,title"
		}

		dataNew, err := utils.GetJSON(dom)

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

		if config.Verbose {
			fmt.Printf("Length of old: %d, new: %d\n", len(Old), len(Posts))
		}
		timeInHours := config.WaitTime.Hours()
		if config.Verbose {
			fmt.Printf("Waittime: %fh\n", timeInHours)
		}
		newPerHour := float64(len(Posts)-len(Old)) / timeInHours
		if config.Verbose {
			fmt.Printf("New posts per hour: %f\n", newPerHour)
		}
		config.StatWpPosts.WithLabelValues(target).Set(newPerHour)
		if config.SaveDiff {
			writeFile(hostname+".diff", newPerHour)
		}

	}
	time.Sleep(config.WaitTime)
}

func RunWPChecksPosts(targets types.URL, config types.WPPostsConfig) {
	go func() {
		for {
			WPPostsFunc(targets, config)

			if config.Wg != nil {
				config.Wg.Done()
			}
		}
	}()
}
