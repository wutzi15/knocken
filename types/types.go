package types

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type URL struct {
	Targets []string `yaml:"targets"`
}

type MetricsConfig struct {
	URLs     URL
	SaveDiff bool
	WaitTime time.Duration
	StatSame *prometheus.GaugeVec
	Verbose  bool
	FastDiff bool
	Wg       *sync.WaitGroup
}

type ContainsConfig struct {
	WaitTime     time.Duration
	Wg           *sync.WaitGroup
	StatContains *prometheus.GaugeVec
	Verbose      bool
}

type KnockenConfig struct {
	Verbose         bool
	SaveDiff        bool
	FastDiff        bool
	WaitTime        time.Duration
	Targets         string
	ContainsTargets string
	Ignore          string
	RunDiff         bool
	RunContain      bool
}

type ContainsTargetSlice []struct {
	Domain  string `yaml:"domain"`
	Contain string `yaml:"contain"`
}

type ContainsTargets struct {
	Targets ContainsTargetSlice `yaml:"targets"`
}
