package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// CF -> for use configs model
	CF = &Configs{}
)

// Environment environment
type Environment string

const (
	// LOCAL environment local
	LOCAL Environment = "local"
	// DEV environment develop
	DEV Environment = "dev"
	// PROD environment production
	PROD Environment = "prod"
)

// Configs models
type Configs struct {
	Minio struct {
		Host       string `mapstructure:"host"`
		AccessKey  string `mapstructure:"access_key"`
		SecretKey  string `mapstructure:"secret_key"`
		BucketName string `mapstructure:"bucket_name"`
		Prefix     string `mapstructure:"prefix"`
		Upload     struct {
			Timestamp string `mapstructure:"timestamp"`
		} `mapstructure:"upload"`
	} `mapstructure:"minio"`
}

// InitConfig init config
func InitConfig(configPath string, environment string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(fmt.Sprintf("config.%s", CF.parseEnvironment(environment)))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingConfig(v, CF); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := bindingConfig(v, CF); err != nil {
			logrus.Error("binding error:", err)
			return
		}
	})

	return nil
}

// bindingConfig binding config
func bindingConfig(vp *viper.Viper, cf *Configs) error {
	if err := vp.Unmarshal(&cf); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}

	return nil
}

func (c Configs) parseEnvironment(environment string) Environment {
	switch environment {
	case "local":
		return LOCAL

	case "dev":
		return DEV

	case "prod":
		return PROD
	}

	return DEV
}
