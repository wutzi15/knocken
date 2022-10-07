package parsers_test

import (
	"strings"
	"testing"

	"github.com/wutzi15/knocken/parsers"
	mytypes "github.com/wutzi15/knocken/types"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func TestTargetParser(t *testing.T) {
	var URLs, err = parsers.ParseTargets("../targets.sample.yml")
	if err != nil {
		t.Errorf("Error parsing targets: %v", err)
	}
	if len(URLs.Targets) != 2 {
		t.Errorf("Expected 2 but got %+v", URLs.Targets)
	}
	if !contains(URLs.Targets, "google.com") {
		t.Errorf("Expected to find google.com but did not %v", URLs.Targets)
	}
	if !contains(URLs.Targets, "escsoftware.de") {
		t.Errorf("Expected to find escsoftware.de but did not %v", URLs.Targets)
	}
}

func TestIgnoreTargets(t *testing.T) {
	URLs, err := parsers.ParseTargets("../targets.sample.yml")
	if err != nil {
		t.Errorf("Error parsing targets: %v", err)
	}

	ignore, err := parsers.ParseTargets("../ignore.sample.yml")
	if err != nil {
		t.Errorf("Error parsing ignore: %v", err)
	}
	URLs = parsers.RemoveIgnoredTargets(URLs, ignore)
	if len(URLs.Targets) != 1 {
		t.Errorf("Expected 1 but got %+v", URLs.Targets)
	}
	if !contains(URLs.Targets, "escsoftware.de") {
		t.Errorf("Expected to find escsoftware.de but did not %v", URLs.Targets)
	}
}

func TestIgnoreFail(t *testing.T) {
	URLs, err := parsers.ParseTargets("foo.yml")
	if len(URLs.Targets) != 0 {
		t.Errorf("Expected 0 but got %+v", URLs.Targets)
	}
	if err == nil {
		t.Errorf("Expected error parsing: %v", err)
	}
}

func containsDomain(d mytypes.ContainsTargetSlice, s string) bool {
	for _, v := range d {
		if strings.Contains(v.Domain, s) {
			return true
		}
	}
	return false
}

func containsContains(d mytypes.ContainsTargetSlice, s string) bool {
	for _, v := range d {
		if strings.Contains(v.Contain, s) {
			return true
		}
	}
	return false
}

func TestContainsTargetParser(t *testing.T) {
	var URLs, err = parsers.ParseContainsTargets("../containsTargets.sample.yml")
	if err != nil {
		t.Errorf("Error parsing targets: %v", err)
	}
	if len(URLs.Targets) != 2 {
		t.Errorf("Expected 2 but got %+v", URLs.Targets)
	}
	if !containsDomain(URLs.Targets, "google.com") {
		t.Errorf("Expected to find google.com but did not %v", URLs.Targets)
	}
	if !containsDomain(URLs.Targets, "escsoftware.de") {
		t.Errorf("Expected to find escsoftware.de but did not %v", URLs.Targets)
	}
	if !containsContains(URLs.Targets, "google") {
		t.Errorf("Expected to find google.com but did not %v", URLs.Targets)
	}
	if !containsContains(URLs.Targets, "Geschäftsführer") {
		t.Errorf("Expected to find escsoftware.de but did not %v", URLs.Targets)
	}
}
