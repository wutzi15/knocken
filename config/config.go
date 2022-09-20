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
	viper.SetDefault("Targets", "targets.yml")
	viper.SetDefault("Ignore", "ignore.yml")
	viper.SetDefault("SaveConfig", false)

	// viper.SetConfigName(".env")
	// viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetEnvPrefix("knocken")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	_ = pflag.Bool("SaveDiff", false, "Keep diffs in ./html/ with diff percentage")
	_ = pflag.Bool("Verbose", false, "Verbose output")
	_ = pflag.String("WaitTime", "5m", "Wait time")
	_ = pflag.String("Targets", "targets.yml", "Targets file")
	_ = pflag.String("Ignore", "ignore.yml", "Ignore file")
	_ = pflag.Bool("SaveConfig", false, "Save config to .env")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("SaveConfig") {
		viper.WriteConfig()
		viper.SafeWriteConfig()
	}
	fmt.Printf("WaitTime: %s\n", viper.GetDuration("WaitTime"))

	return types.KnockenConfig{
		Verbose:  viper.GetBool("Verbose"),
		SaveDiff: viper.GetBool("SaveDiff"),
		WaitTime: viper.GetDuration("WaitTime"),
		Targets:  viper.GetString("Targets"),
		Ignore:   viper.GetString("Ignore"),
	}

}
