package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var config *Config
var once sync.Once

// Get reads the config.yml, overrides values with environment variables
// and returns a singleton Config struct
func Get() *Config {
	once.Do(func() {
		addConfigSearchPath()
		setViperConfig()
		setDefaultValues()

		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		viper.AutomaticEnv()

		config = &Config{}
		err = viper.Unmarshal(config)

		if err != nil {
			panic(err)
		}
		if config.Log.Level == "" {
			config.Log.Level = "info"
		}
	})
	return config
}

func addConfigSearchPath() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
}

func setViperConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func setDefaultValues() {
	viper.SetDefault("APPLICATION_NAME", "elastic-talk-search")
	viper.SetDefault("HOSTNAME", "not_set")
}
