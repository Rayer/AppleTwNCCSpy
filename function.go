package AppleProductMonitor

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

func init() {
	functions.CloudEvent("AppleProductMonitor", CrawlAndAnalyze)
}

//go:embed resources/newitem_explanation.txt
var NewItemExplanation string

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Configuration struct {
	BotToken         string `yaml:"botToken"`
	ChannelId        int64  `yaml:"channelId"`
	FetchTarget      string `yaml:"fetchTarget"`
	FetchIntervalSec int    `yaml:"fetchIntervalSec"`
	DebugLevel       int    `yaml:"debugLevel"`
	Bucket           string `yaml:"bucket"`
	Prefix           string `yaml:"prefix"`
}

func CrawlAndAnalyze(ctx context.Context, e event.Event) error {
	log.Printf("received nccEvents %+v\n", e)
	confByte, err := os.ReadFile("/var/secrets/ncc-bot-config/config.yml")
	if err != nil {
		// List the directory
		files, _ := os.ReadDir("/var/secrets/ncc-bot-config")
		return errors.New(fmt.Sprintf("config file not found, files = %+v", files))
	}

	log.Printf("config : %s", string(confByte))
	config := Configuration{}
	err = yaml.Unmarshal(confByte, &config)
	if err != nil {
		return err
	}

	//Initialize bot
	bot, err := tgbotapi.NewBotAPI(config.BotToken)

	if err != nil {
		return fmt.Errorf("bot initialization failed, bot ID is %s, err = %v", config.BotToken, err)
	}

	//Initialize Storage
	dataAccess, err := NewGcsDataAccess(ctx, config.Bucket, config.Prefix)
	if err != nil {
		return fmt.Errorf("GCS initialization failed, err = %v", err)
	}

	crawler := Crawler{
		DataAccess:  dataAccess,
		FetchTarget: "https://www.apple.com/tw/nccid",
	}

	nccEvents, err := crawler.FetchAndCompare(ctx)

	if err != nil {
		return err
	}

	if len(nccEvents.Added) > 0 {
		sendToChannel(bot, config.ChannelId, "偵測到新加入產品:")
		addList := prettyPrintProducts(nccEvents.Added, "(+)")
		for _, v := range addList {
			sendToChannel(bot, config.ChannelId, v)
		}
		sendToChannel(bot, config.ChannelId, NewItemExplanation)
	}
	if len(nccEvents.Removed) > 0 {
		sendToChannel(bot, config.ChannelId, "偵測到移除產品:")
		removeList := prettyPrintProducts(nccEvents.Removed, "(-)")
		for _, v := range removeList {
			sendToChannel(bot, config.ChannelId, v)
		}
	}

	if len(nccEvents.Removed) != 0 || len(nccEvents.Added) != 0 {
		//Serialize into gcs
		err = dataAccess.SaveDiff(ctx, nccEvents)
		if err != nil {
			return fmt.Errorf("unable write diff file : %v", err)
		}
	}

	return nil
}

func sendToChannel(bot *tgbotapi.BotAPI, channelId int64, message string) {
	log.Println("Sending : ", message)
	msg := tgbotapi.NewMessage(channelId, message)
	_, _ = bot.Send(msg)
	time.Sleep(1 * time.Second)
}

func prettyPrintProducts(source []Product, prefix string) (ret []string) {
	log.Printf("Trying to print : %+v", source)
	productMap := make(map[string][]Product)
	for _, r := range source {
		if v, exist := productMap[r.Group]; !exist {
			productMap[r.Group] = []Product{r}
		} else {
			productMap[r.Group] = append(v, r)
		}
	}

	log.Printf("Product Map : %+v", productMap)

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
