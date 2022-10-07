package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wutzi15/knocken/types"
)

func GetConfig() types.KnockenConfig {
	viper.SetDefault("Verbose", false)
	viper.SetDefault("SaveDiff", false)
	viper.SetDefault("FastDiff", false)
	viper.SetDefault("WaitTime", "5m")
	viper.SetDefault("Targets", "targets.yml")
	viper.SetDefault("ContainsTargets", "containstargets.yml")
	viper.SetDefault("Ignore", "ignore.yml")
	viper.SetDefault("SaveConfig", false)
	viper.SetDefault("RunDiff", true)
	viper.SetDefault("RunContain", true)
	viper.SetDefault("RunWPPosts", true)

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetEnvPrefix("knocken")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	flag := pflag.FlagSet{}

	_ = flag.Bool("SaveDiff", false, "Keep diffs in ./html/ with diff percentage")
	_ = flag.Bool("FastDiff", false, "Use faster diff algorithm (Jaro-Winkler )")
	_ = flag.Bool("Verbose", false, "Verbose output")
	_ = flag.String("WaitTime", "5m", "Wait time")
	_ = flag.String("Targets", "targets.yml", "Targets file")
	_ = flag.String("ContainsTargets", "containstargets.yml", "Targets file for the contains check")
	_ = flag.String("Ignore", "ignore.yml", "Ignore file")
	_ = flag.Bool("SaveConfig", false, "Save config to .env")
	_ = flag.Bool("RunDiff", true, "Run the diff check")
	_ = flag.Bool("RunContain", true, "Run the contains check")
	_ = flag.Bool("RunWPPosts", true, "Run the wordpress posts check")

	flag.Parse(os.Args[1:])
	viper.BindPFlags(&flag)

	if viper.GetBool("SaveConfig") {
		viper.WriteConfig()
		viper.SafeWriteConfig()
	}

	config := types.KnockenConfig{
		Verbose:         viper.GetBool("Verbose"),
		SaveDiff:        viper.GetBool("SaveDiff"),
		FastDiff:        viper.GetBool("FastDiff"),
		WaitTime:        viper.GetDuration("WaitTime"),
		Targets:         viper.GetString("Targets"),
		ContainsTargets: viper.GetString("ContainsTargets"),
		Ignore:          viper.GetString("Ignore"),
		RunDiff:         viper.GetBool("RunDiff"),
		RunContain:      viper.GetBool("RunContain"),
		RunWPPosts:      viper.GetBool("RunWPPosts"),
	}

	if viper.GetBool("Verbose") {
		fmt.Println("Verbose output enabled")
		fmt.Printf("Config: %+v\n", config)
	}

	return config

}
