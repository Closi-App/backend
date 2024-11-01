package config

import "github.com/spf13/viper"

func NewConfig(cfgFilePath string) *viper.Viper {
	cfg := viper.New()
	cfg.SetConfigFile(cfgFilePath)
	if err := cfg.ReadInConfig(); err != nil {
		panic("error reading config: " + err.Error())
	}

	return cfg
}
