package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wutzi15/knocken/types"
)

func GetConfig() types.KnockenConfig {
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

	_ = pflag.Bool("SaveDiff", false, "Keep diffs in ./html/ with diff percentage")
	_ = pflag.Bool("Verbose", false, "Verbose output")
	_ = pflag.String("WaitTime", "5m", "Wait time")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	fmt.Printf("WaitTime: %s\n", viper.Get("WaitTime"))

	return types.KnockenConfig{
		Verbose:  viper.GetBool("Verbose"),
		SaveDiff: viper.GetBool("SaveDiff"),
		WaitTime: viper.GetDuration("WaitTime"),
	}

}
