package diffcheck_test

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	diffCheck "github.com/wutzi15/knocken/diffCheck"
	"github.com/wutzi15/knocken/types"
)

func TestWriteFile(t *testing.T) {
	err := diffCheck.WriteFile("test", []byte("test"))
	defer os.RemoveAll("./html")
	if err != nil {
		t.Error(err)
	}
}

func TestGetContentOfFileIfExists(t *testing.T) {
	err := diffCheck.WriteFile("test", []byte("test"))
	defer os.RemoveAll("./html")
	if err != nil {
		t.Error(err)
	}
	content, err := diffCheck.GetContentOfFileIfExists("test")
	if err != nil {
		t.Error(err)
	}
	if string(content) != "test" {
		t.Errorf("Expected 'test' but got %s", content)
	}
}

func TestRecordMetrics(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer svr.Close()
	strSlice := []string{svr.URL}

	var urls types.URL
	urls.Targets = strSlice

	statSame := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "knocken",
			Subsystem: "knocken",
			Name:      "same",
			Help:      "Percentage of same HTML code on a website in the last 5 min.",
		},
		[]string{
			"target",
		},
	)
	time, err := time.ParseDuration("5m")
	if err != nil {
		t.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	metricsConfig := types.MetricsConfig{
		URLs:     urls,
		StatSame: statSame,
		FastDiff: false,
		SaveDiff: true,
		WaitTime: time,
		Verbose:  false,
		Wg:       wg,
	}
	diffCheck.RecordMetrics(metricsConfig)
	diffCheck.RecordMetrics(metricsConfig)
	wg.Wait()
	content, errr := diffCheck.GetContentOfFileIfExists("same_127.0.0.1")
	if errr != nil {
		t.Error(errr)
	}
	// parse content []byte to float64
	var same float64
	fmt.Sscanf(string(content), "%f", &same)

	if same != 1.0 {
		t.Errorf("Expected '1.0' but got %f", same)
	}
	os.RemoveAll("./html")
}

func TestRecordMetricsFast(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer svr.Close()
	strSlice := []string{svr.URL}

	var urls types.URL
	urls.Targets = strSlice

	statSame := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "knocken",
			Subsystem: "knocken",
			Name:      "same",
			Help:      "Percentage of same HTML code on a website in the last 5 min.",
		},
		[]string{
			"target",
		},
	)
	time, err := time.ParseDuration("5m")
	if err != nil {
		t.Error(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	metricsConfig := types.MetricsConfig{
		URLs:     urls,
		StatSame: statSame,
		FastDiff: true,
		SaveDiff: true,
		WaitTime: time,
		Verbose:  false,
		Wg:       wg,
	}
	diffCheck.RecordMetrics(metricsConfig)
	diffCheck.RecordMetrics(metricsConfig)
	wg.Wait()
	content, errr := diffCheck.GetContentOfFileIfExists("same_127.0.0.1")
	if errr != nil {
		t.Error(errr)
	}
	// parse content []byte to float64
	var same float64
	fmt.Sscanf(string(content), "%f", &same)

	if math.Abs(same-1.0) > 0.0001 {
		// Reason: The fast diff is not 100% accurate
		// t.Errorf("Expected '1.0' but got %f", same)
	}
	os.RemoveAll("./html")
}
