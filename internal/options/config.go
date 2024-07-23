package options

import (
	"strings"

	"github.com/spf13/viper"
)

type ConfigYaml struct {
	Core BindConfigRest
}

type BindConfigRest struct {
	BindAddress string `yaml:"bindaddr"`
	BindPort    int    `yaml:"bindport"`
	MetricsPath string `yaml:"metricspath"`
}

func LoadConfig() (*ConfigYaml, error) {
	conf := &ConfigYaml{}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("smartwatteco")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("/")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetDefault("cire.bindaddr", "0.0.0.0")
	viper.SetDefault("core.bindport", 8080)
	viper.SetDefault("core.metricspath", "/metrics")

	conf.Core.BindAddress = viper.GetString("core.bindaddr")
	conf.Core.BindPort = viper.GetInt("core.bindport")
	conf.Core.MetricsPath = viper.GetString("core.metricspath")

	return conf, nil
}
