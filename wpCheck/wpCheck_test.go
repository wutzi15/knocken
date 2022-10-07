package wpcheck_test

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/wutzi15/knocken/types"
	wpcheck "github.com/wutzi15/knocken/wpCheck"
)

func TestWPLong(t *testing.T) {

	strSlice := []string{"https://www.escsoftware.de"}

	var urls types.URL
	urls.Targets = strSlice

	statwpPosts := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "knocken",
			Subsystem: "knocken",
			Name:      "wpPosts",
			Help:      "number of new posts per 1 hour on a wordpress website. ",
		},
		[]string{
			"target",
		},
	)

	time, err := time.ParseDuration("5s")
	if err != nil {
		t.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	wpcfg := types.WPPostsConfig{
		WaitTime:    time,
		Verbose:     false,
		StatWpPosts: statwpPosts,
		Wg:          wg,
		SaveDiff:    true,
		Testing:     true,
	}
	wpcheck.RunWPChecksPosts(urls, wpcfg)
	wpcheck.RunWPChecksPosts(urls, wpcfg)

	// os.RemoveAll("./json")
}
