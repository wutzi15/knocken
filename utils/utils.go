package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
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

func GetJSON(url string) ([]byte, error) {
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
