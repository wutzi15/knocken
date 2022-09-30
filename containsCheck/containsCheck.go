package containscheck

import (
	"fmt"
	"strings"
	"time"

	"github.com/wutzi15/knocken/types"
	"github.com/wutzi15/knocken/utils"
)

func ContainsFunc(targets types.ContainsTargets, config types.ContainsConfig) []bool {
	var ret []bool
	for _, target := range targets.Targets {
		html, err := utils.GetHTML(target.Domain)
		if err != nil {
			fmt.Printf("Error getting html from %s: %s)", target.Domain, err)
		}
		contains := strings.Contains(string(html), target.Contain)
		var cnt float64 = 0.0
		if contains {
			cnt = 1.0
		}
		config.StatContains.WithLabelValues(target.Domain).Set(cnt)
		ret = append(ret, contains)
	}
	return ret
}

func RunContain(targets types.ContainsTargets, config types.ContainsConfig) {
	go func() {
		for {
			_ = ContainsFunc(targets, config)

			if config.Wg != nil {
				config.Wg.Done()
			}
			time.Sleep(config.WaitTime)
		}
	}()
}
