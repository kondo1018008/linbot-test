package main

import (
	"encoding/json"
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
	errCheck(err)
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
				errCheck(err)

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

	replyMsg := getRestoInfo(lat,lng)
	_, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage("以下が近辺でお勧めの飲食店です")).Do()
	errCheck(err)

	_, err = bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
	errCheck(err)
}

func getRestoInfo(lat string, lng string) string{
	key := os.Getenv("API_KEY")
	url := fmt.Sprintf("http://webservice.recruit.co.jp/hotpepper/gourmet/v1/?key=%s&lat=%s&lng=%s&range=5&order=4&format=json",key, lat, lng)
	fmt.Println(url)

	resp, err := http.Get(url)
	errCheck(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	errCheck(err)

	var data response
	err = json.Unmarshal(body,&data)
	errCheck(err)

	info := ""
	for _, shop := range data.Results.Shop {
		info += shop.Name + "\n" + shop.Address + "\n\n"
	}

	return info
}

func errCheck(err error){
	if err != nil {
		log.Fatal(err)
	}
}

type response struct {
	Results results `json:"results"`
}

type results struct {
	Shop []shop `json:"shop"`
}

type shop struct {
	Name string `json:"name"`
	Address string `json:"address"`
}