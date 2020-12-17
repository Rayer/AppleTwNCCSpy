package main

import (
	"AppleProductMonitor"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

//Real one
//const chatId = int64(-1001210598225)

//Test channel
//const chatId = int64(-1001409488439)

func main() {
	//log.SetLevel(log.DebugLevel)
	conf, err := NewConfigurationFromViper()

	if err != nil {
		panic(err)
	}

	log.Infof("Loaded configuration : %+v", conf)
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)

	if err != nil {
		log.Panic(err)
	}

	log.SetLevel(log.Level(conf.DebugLevel))
	bot.Debug = true
	log.Printf("Successfully authorized on account %s", bot.Self.UserName)

	c := AppleProductMonitor.NewCrawler(conf.FetchTarget, conf.FetchIntervalSec)
	ch := c.Run(context.TODO())

	for {
		select {
		case event := <-ch:
			if len(event.Added) > 0 {
				sendToChannel(bot, conf.ChannelId, "偵測到新加入產品:")
				addList := prettyPrintProducts(event.Added, "(+)")
				for _, v := range addList {
					sendToChannel(bot, conf.ChannelId, v)
				}
			}
			if len(event.Removed) > 0 {
				sendToChannel(bot, conf.ChannelId, "偵測到移除產品:")
				removeList := prettyPrintProducts(event.Removed, "(-)")
				for _, v := range removeList {
					sendToChannel(bot, conf.ChannelId, v)
				}
			}
		}
	}
}

func sendToChannel(bot *tgbotapi.BotAPI, chatId int64, message string) {
	log.Info("Sending : ", message)
	msg := tgbotapi.NewMessage(chatId, message)
	_, _ = bot.Send(msg)
	time.Sleep(1 * time.Second)
}

func prettyPrintProducts(source []AppleProductMonitor.Product, prefix string) (ret []string) {
	log.Infof("Trying to print : %+v", source)
	productMap := make(map[string][]AppleProductMonitor.Product)
	for _, r := range source {
		if v, exist := productMap[r.Group]; !exist {
			productMap[r.Group] = []AppleProductMonitor.Product{r}
		} else {
			productMap[r.Group] = append(v, r)
		}
	}

	log.Infof("Product Map : %+v", productMap)

	ret = make([]string, 0)

	for k, v := range productMap {
		var buffer string
		buffer += fmt.Sprintf("產品 : %s\n", k)
		for _, p := range v {
			buffer += fmt.Sprintf("%s%s, %s, %s\n", prefix, p.Model, p.NCC, p.Product)
		}
		ret = append(ret, buffer)
	}
	return
}
