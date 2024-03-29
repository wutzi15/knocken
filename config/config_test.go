package config_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wutzi15/knocken/config"
)

func TestGetConfig(t *testing.T) {
	config := config.GetConfig()
	if config.Verbose != false {
		t.Errorf("Expected false but got %v", config.Verbose)
	}
	if config.SaveDiff != false {
		t.Errorf("Expected false but got %v", config.SaveDiff)
	}
	dur, _ := time.ParseDuration("5m")
	if config.WaitTime != dur {
		t.Errorf("Expected 5m but got %v", config.WaitTime)
	}
	if config.Targets != "targets.yml" {
		t.Errorf("Expected 'targets.yml' but got %v", config.Targets)
	}
	if config.ContainsTargets != "containstargets.yml" {
		t.Errorf("Expected 'constainstargets.yml' but got %v", config.Targets)
	}
	if config.Ignore != "ignore.yml" {
		t.Errorf("Expected 'ignore.yml' but got %v", config.Ignore)
	}
	if config.RunDiff != true {
		t.Errorf("Expected 'true' but got %v", config.RunDiff)
	}
	if config.RunContain != true {
		t.Errorf("Expected 'true' but got %v", config.RunContain)
	}
	if config.RunWPPosts != true {
		t.Errorf("Expected 'true' but got %v", config.RunWPPosts)
	}
	if config.WPTargets != "wptargets.yml" {
		t.Errorf("Expected 'true' but got %v", config.WPTargets)
	}
}

func TestConfigFromEnv(t *testing.T) {
	os.Setenv("KNOCKEN_VERBOSE", "true")
	os.Setenv("KNOCKEN_SAVEDIFF", "true")
	os.Setenv("KNOCKEN_FASTDIFF", "true")
	os.Setenv("KNOCKEN_WAITTIME", "7m")
	os.Setenv("KNOCKEN_TARGETS", "foo")
	os.Setenv("KNOCKEN_CONTAINSTARGETS", "baz")
	os.Setenv("KNOCKEN_IGNORE", "bar")
	os.Setenv("KNOCKEN_RUNDIFF", "false")
	os.Setenv("KNOCKEN_RUNCONTAIN", "false")
	os.Setenv("KNOCKEN_RUNWPPOSTS", "false")
	os.Setenv("KNOCKEN_WPTARGETS", "nupf")
	defer func() {
		os.Unsetenv("KNOCKEN_VERBOSE")
		os.Unsetenv("KNOCKEN_SAVEDIFF")
		os.Unsetenv("KNOCKEN_FASTDIFF")
		os.Unsetenv("KNOCKEN_WAITTIME")
		os.Unsetenv("KNOCKEN_TARGETS")
		os.Unsetenv("KNOCKEN_CONTAINSTARGETS")
		os.Unsetenv("KNOCKEN_IGNORE")
		os.Unsetenv("KNOCKEN_RUNDIFF")
		os.Unsetenv("KNOCKEN_RUNCONTAIN")
		os.Unsetenv("KNOCKEN_RUNWPPOSTS")
		os.Unsetenv("KNOCKEN_WPTARGETS")
	}()
	config := config.GetConfig()
	if config.Verbose != true {
		t.Errorf("Expected false but got %v", config.Verbose)
	}
	if config.SaveDiff != true {
		t.Errorf("Expected false but got %v", config.SaveDiff)
	}
	if config.FastDiff != true {
		t.Errorf("Expected false but got %v", config.FastDiff)
	}
	dur, _ := time.ParseDuration("7m")
	if config.WaitTime != dur {
		t.Errorf("Expected 7m but got %v", config.WaitTime)
	}
	if config.Targets != "foo" {
		t.Errorf("Expected 'foo' but got %v", config.Targets)
	}
	if config.ContainsTargets != "baz" {
		t.Errorf("Expected 'foo' but got %v", config.Targets)
	}
	if config.Ignore != "bar" {
		t.Errorf("Expected 'bar' but got %v", config.Ignore)
	}
	if config.RunDiff != false {
		t.Errorf("Expected 'false' but got %v", config.RunDiff)
	}
	if config.RunContain != false {
		t.Errorf("Expected 'false' but got %v", config.RunContain)
	}
	if config.RunWPPosts != false {
		t.Errorf("Expected 'false' but got %v", config.RunWPPosts)
	}
	if config.WPTargets != "nupf" {
		t.Errorf("Expected 'false' but got %v", config.WPTargets)
	}
}

func TestConfigWrite(t *testing.T) {
	os.Remove(".env")
	defer os.Remove(".env")
	os.Setenv("KNOCKEN_SAVECONFIG", "true")
	config.GetConfig()
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		t.Errorf("Expected .env to exist")
	}
	data, err := os.ReadFile(".env")
	if err != nil {
		t.Errorf("Error reading .env")
	}
	expect, err := os.ReadFile("../env.sample")
	if err != nil {
		t.Errorf("Error reading env.sample")
	}
	if strings.TrimSpace(string(data)) != strings.TrimSpace(string(expect)) {
		t.Errorf("Expected %s but got %s", expect, data)
	}
}

func TestConfigReadEnv(t *testing.T) {
	os.Remove(".env")
	defer func() {
		os.Remove(".env")
	}()
	out := `
		# knocken config file
		IGNORE=baz
		SAVECONFIG=true
		SAVEDIFF=false
		FASTDIFF=true
		TARGETS=nupf
		CONTAINSTARGETS=foo
		VERBOSE=false
		WAITTIME=7m
		RUNDIFF=false
		RUNCONTAIN=false
		RUNWPPOSTS=false
		WPTARGETS=bar
	`
	os.WriteFile(".env", []byte(out), 0644)
	config := config.GetConfig()
	if config.Verbose != false {
		t.Errorf("Expected false but got %v", config.Verbose)
	}
	if config.SaveDiff != false {
		t.Errorf("Expected false but got %v", config.SaveDiff)
	}
	if config.FastDiff != true {
		t.Errorf("Expected false but got %v", config.FastDiff)
	}
	dur, _ := time.ParseDuration("7m")
	if config.WaitTime != dur {
		t.Errorf("Expected 7m but got %v", config.WaitTime)
	}
	if config.Targets != "nupf" {
		t.Errorf("Expected 'nupf' but got %v", config.Targets)
	}
	if config.ContainsTargets != "foo" {
		t.Errorf("Expected 'foo' but got %v", config.ContainsTargets)
	}
	if config.Ignore != "baz" {
		t.Errorf("Expected 'baz' but got %v", config.Ignore)
	}
	if config.RunDiff != false {
		t.Errorf("Expected 'false' but got %v", config.RunDiff)
	}
	if config.RunContain != false {
		t.Errorf("Expected 'false' but got %v", config.RunContain)
	}
	if config.RunWPPosts != false {
		t.Errorf("Expected 'false' but got %v", config.RunWPPosts)
	}
	if config.WPTargets != "bar" {
		t.Errorf("Expected 'bar' but got %v", config.WPTargets)
	}
}
