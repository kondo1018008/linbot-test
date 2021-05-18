package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main(){
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/callback",lineHandler)

	fmt.Println("running in http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080",nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request){
	msg := "Hello world"
	fmt.Fprintln(w, msg)
}

func lineHandler(w http.ResponseWriter, r *http.Request){
	secret := os.Getenv("CHANNEL_SECRET")
	token := os.Getenv("CHANNEL_TOKEN")
	bot, err := linebot.New(
		secret,
		token,
		)
	if err != nil {
		log.Fatal(err)
	}
	// リクエストからBOTのイベントを取得
	events, err := bot.ParseRequest(r)
	// リクエストのチェック
	if err != nil {
		if err == linebot.ErrInvalidSignature{
			w.WriteHeader(400)
		}else{
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type){
			case *linebot.TextMessage:
				replyMessage := message.Text
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Fatal(err)
				}
			case *linebot.LocationMessage:
				sendRestoInfo(bot, event)
			}
		}
	}
}

func sendRestoInfo(bot *linebot.Client, e *linebot.Event){
	msg := e.Message.(*linebot.LocationMessage)

	lat := strconv.FormatFloat(msg.Latitude, 'f', 2, 64)
	lng := strconv.FormatFloat(msg.Longitude, 'f', 2, 64)

	replyMsg := fmt.Sprintf("経度：%s\n緯度：%s", lat, lng)

	_, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
	if err != nil {
		log.Println(err)
	}
	key := os.Getenv("API_KEY")
	url := fmt.Sprintf("http://webservice.recruit.co.jp/hotpepper/gourmet/v1/?key=%s&lat=%s&lng=%s&range=5&order=4&count=1",key, lat, lng)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	restrans, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(restrans))


}