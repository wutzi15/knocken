package containscheck_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	containscheck "github.com/wutzi15/knocken/containsCheck"
	"github.com/wutzi15/knocken/types"
)

func TestContainsCheck(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}))
	defer svr.Close()

	tgts := types.ContainsTargetSlice{
		{
			svr.URL,
			"test",
		},
	}

	var urls types.ContainsTargets = types.ContainsTargets{
		tgts,
	}

	containsSame := prometheus.NewGaugeVec(
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

	cfg := types.ContainsConfig{
		StatContains: containsSame,
		WaitTime:     time,
		Verbose:      false,
		Wg:           nil,
	}

	cnt := containscheck.ContainsFunc(urls, cfg)
	if len(cnt) != 1 {
		t.Error("Expected 1 result, got ", len(cnt))
	}

	if cnt[0] != true {
		t.Error("Expected true, got ", cnt[0])
	}
}
