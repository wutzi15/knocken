package types

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type URL []struct {
	Targets []string `yaml:"targets"`
}

type MetricsConfig struct {
	URLs     URL
	SaveDiff bool
	WaitTime time.Duration
	StatSame *prometheus.GaugeVec
	Verbose  bool
	Wg       *sync.WaitGroup
}

type KnockenConfig struct {
	Verbose  bool
	SaveDiff bool
	WaitTime time.Duration
}
