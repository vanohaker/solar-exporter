package options

import (
	"os"
	"strings"

	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

type ConfigYaml struct {
	Core BindConfigRest
}

type BindConfigRest struct {
	BindAddress   string `yaml:"bindaddr"`
	BindPort      int    `yaml:"bindport"`
	MetricsPath   string `yaml:"metricspath"`
	MetricsPrefix string `yaml:"metricsprefix"`
}

func LoadConfig() (*ConfigYaml, error) {
	conf := &ConfigYaml{}

	envprefix := os.Getenv("ENVPREFIX")
	if envprefix == "" {
		envprefix = "solarexporter"
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envprefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("/")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		fiberlog.Info(err.Error())
	}

	viper.SetDefault("cire.bindaddr", "0.0.0.0")
	viper.SetDefault("core.bindport", 9560)
	viper.SetDefault("core.metricspath", "/metrics")
	viper.SetDefault("core.metricsprefix", "solar_invertor")

	conf.Core.BindAddress = viper.GetString("core.bindaddr")
	conf.Core.BindPort = viper.GetInt("core.bindport")
	conf.Core.MetricsPath = viper.GetString("core.metricspath")
	conf.Core.MetricsPrefix = viper.GetString("core.metricsprefix")

	return conf, nil
}
