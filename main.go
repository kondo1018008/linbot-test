package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
)

func main(){
	http.HandleFunc("/callback",lineHandler)

	fmt.Println("running in http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080",nil))
}

func lineHandler(w http.ResponseWriter, r *http.Response){
	secret := os.Getenv("CHANNEL_SECRET")
	token := os.Getenv("CHANNEL_TOKEN")
	bot, err := linebot.New(
		secret,
		token,
		)
	if err != nil {
		log.Fatal(err)
	}

}
