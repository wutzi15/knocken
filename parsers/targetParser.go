package parsers

import (
	"io/ioutil"

	"github.com/wutzi15/knocken/types"
	"gopkg.in/yaml.v2"
)

func ParseTargets(targets string) (types.URL, error) {
	data, err := ioutil.ReadFile(targets)
	if err != nil {
		return types.URL{}, err
	}

	var URLs types.URL
	err = yaml.Unmarshal(data, &URLs)
	if err != nil {
		return types.URL{}, err
	}
	return URLs, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func RemoveIgnoredTargets(URLs types.URL, ignore types.URL) types.URL {
	var newURLs types.URL
	for i := 0; i < len(URLs.Targets); i++ {
		if !contains(ignore.Targets, URLs.Targets[i]) {
			newURLs.Targets = append(newURLs.Targets, URLs.Targets[i])
		}
	}

	return newURLs
}
