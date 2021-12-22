package main

import (
	"context"
	_ "embed"
	"fmt"
	AppleProductMonitor "github.com/Rayer/AppleTwNCCSpy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

//Real one
//const chatId = int64(-1001210598225)

//Test channel
//const chatId = int64(-1001409488439)

//go:embed resources/newitem_explanation.txt
var NewItemExplanation string

type BotSpec struct {
	config        *Configuration
	bot           *tgbotapi.BotAPI
	context       context.Context
	cancel        context.CancelFunc
	crawler       *AppleProductMonitor.Crawler
	updateChannel chan AppleProductMonitor.Event
}

func (b *BotSpec) teardown() {
	if b.cancel != nil {
		b.cancel()
	}
}

func NewBotSpec(conf *Configuration) *BotSpec {

	b := &BotSpec{
		config: conf,
	}

	var err error

	b.bot, err = tgbotapi.NewBotAPI(conf.BotToken)

	if err != nil {
		log.Panic(err)
	}

	log.SetLevel(log.Level(conf.DebugLevel))
	b.bot.Debug = true
	log.Infof("Successfully authorized on account %s", b.bot.Self.UserName)

	b.crawler = AppleProductMonitor.NewCrawler(conf.FetchTarget, conf.FetchIntervalSec)
	b.context, b.cancel = context.WithCancel(context.Background())
	b.updateChannel = b.crawler.Run(context.Background())

	return b
}

var botSpec *BotSpec

func main() {
	conf, err := NewConfigurationFromViper(func(configuration *Configuration) {
		botSpec.teardown()
		log.Infof("New configuration : %+v", configuration)
		botSpec = NewBotSpec(configuration)
	})

	if err != nil {
		panic(err)
	}

	log.Infof("Loaded configuration : %+v", conf)

	botSpec = NewBotSpec(conf)

	for {
		select {
		case event := <-botSpec.updateChannel:
			if len(event.Added) > 0 {
				botSpec.sendToChannel("偵測到新加入產品:")
				addList := botSpec.prettyPrintProducts(event.Added, "(+)")
				for _, v := range addList {
					botSpec.sendToChannel(v)
				}
				botSpec.sendToChannel(NewItemExplanation)
			}
			if len(event.Removed) > 0 {
				botSpec.sendToChannel("偵測到移除產品:")
				removeList := botSpec.prettyPrintProducts(event.Removed, "(-)")
				for _, v := range removeList {
					botSpec.sendToChannel(v)
				}
			}
		default:
		}
	}
}

func (b *BotSpec) sendToChannel(message string) {
	log.Info("Sending : ", message)
	msg := tgbotapi.NewMessage(b.config.ChannelId, message)
	_, _ = b.bot.Send(msg)
	time.Sleep(1 * time.Second)
}

func (b *BotSpec) prettyPrintProducts(source []AppleProductMonitor.Product, prefix string) (ret []string) {
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
