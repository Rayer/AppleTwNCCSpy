# Apple NCC Product Monitor

一個會自動爬[Apple NCC產品清單](https://www.apple.com/tw/nccid) 的爬蟲。以Model(如A2421)為主，如果有偵測到任何添加或者減少，就會以telegram bot告知特定頻道。

## 使用方法

### 首先，建立Telegram Bot跟頻道

請參考[這裏](https://dotblogs.com.tw/c5todo/2016/07/10/174049)以取得你的Telegram Bot，建立頻道，以及取得chat id。

取得chat id有簡單一點的方法就是，傳訊給@getidsbot，然後將你想要取得的channel/group裡面任何一條訊息forward給他，他就會告訴你該channel/group的詳細訊息，包含chat id。

### 組態檔案

固定放置在 ./valut或者工作目錄底下。組態檔固定檔名為bot.yaml，內容如下

```
BotToken: "<你的Bot token>"
ChannelId: <Channel ID>
DebugLevel: 1
```

範例 (都假的）:

```
BotToken: "123456789:AABCDEFUGA-Z83AjgisIIS"
ChannelId: -12341234111
DebugLevel: 1
```

### Native版本

直接編譯bot/bot_main.go即可，把組態檔放入工作目錄下。在工作目錄底下建立/vault資料夾讓他可以寫資料進去即可。

### Docker版本

由於這個tag我自己也在用，我不保證相容性，所以還是建議自己build docker image。

```
# tag可以自己取名，或者直接docker pull rayer/apple-product-monitor也可以
docker build . -t rayer/apple-product-monitor
# 請找一個目錄bind mount進/app/vault，如底下-v參數。把設定檔放入該目錄即可。
# 以下面指令為例子的話 就是建立/opt/apple-product-monitor，把bot.yaml放入裡面
docker run --name AppleProductMonitor -v /opt/apple-product-monitor:/app/vault --rm -d rayer/apple-product-monitor
```

### Google Cloud Functions版本

這個比較牽涉到設定，不過可以參考 `function.go` 這隻。我自己的做法是使用Secret Manager

1. 參考 `function.go` 裡面的 `Configuration` 寫出config的yaml內容，寫入Secret Manager一個值裡面
2. 將這個內容mount在 `/secrets/AppleTwNccBotConfig.yml` 具體作法請參照Cloud Function Config
3. 參照Config的 `bucket` `prefix` 在GCS上建立bucket跟folder

通常這樣就可以了，剩下的要自己跑跑看。我自己是跑一套在GCP上啦，問題不大。

