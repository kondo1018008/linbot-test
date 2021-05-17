package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
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
			}
		}
	}


}
