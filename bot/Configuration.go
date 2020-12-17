package main

import "github.com/spf13/viper"

type Configuration struct {
	BotToken         string
	ChannelId        int64
	FetchTarget      string
	FetchIntervalSec int
	DebugLevel       int
}

func NewConfigurationFromViper() (conf *Configuration, err error) {
	viper.New()
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("./vault")
	viper.SetConfigName("bot")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	viper.SetDefault("FetchTarget", "https://www.apple.com/tw/nccid")
	viper.SetDefault("FetchIntervalSec", 600)

	conf = &Configuration{
		BotToken:         viper.GetString("BotToken"),
		ChannelId:        viper.GetInt64("ChannelId"),
		FetchTarget:      viper.GetString("FetchTarget"),
		FetchIntervalSec: viper.GetInt("FetchIntervalSec"),
		DebugLevel:       viper.GetInt("DebugLevel"),
	}

	return
}
