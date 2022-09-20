package config

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
saveDiff := flag.Bool("saveDiffs", false, "Keep diffs in ./html/ with diff percentage")
	v := flag.Bool("verbose", false, "Verbose output")
	waitTimeStr := flag.String("waitTime", "5m", "Wait time")
*/

func SetUpConfig() {
	viper.SetDefault("Verbose", false)
	viper.SetDefault("SaveDiff", false)
	viper.SetDefault("WaitTime", "5m")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	viper.SetEnvPrefix("knocken")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// viper.WriteConfig()
	// viper.SafeWriteConfig()
}
