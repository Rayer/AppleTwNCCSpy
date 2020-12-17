package main

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	BotToken         string
	ChannelId        int64
	FetchTarget      string
	FetchIntervalSec int
	DebugLevel       int
}

func NewConfigurationFromViper(onConfigChanged func(configuration *Configuration)) (conf *Configuration, err error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("./vault")
	v.SetConfigName("bot")
	v.SetConfigType("yaml")

	err = v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	conf = createConfigurationFromViper(v)
	if onConfigChanged != nil {
		v.WatchConfig()
		v.OnConfigChange(func(in fsnotify.Event) {
			log.Info("Configuration changed detected...")
			onConfigChanged(createConfigurationFromViper(v))
		})
	}

	return
}

func createConfigurationFromViper(v *viper.Viper) *Configuration {

	v.SetDefault("FetchTarget", "https://www.apple.com/tw/nccid")
	v.SetDefault("FetchIntervalSec", 600)

	return &Configuration{
		BotToken:         v.GetString("BotToken"),
		ChannelId:        v.GetInt64("ChannelId"),
		FetchTarget:      v.GetString("FetchTarget"),
		FetchIntervalSec: v.GetInt("FetchIntervalSec"),
		DebugLevel:       v.GetInt("DebugLevel"),
	}
}
