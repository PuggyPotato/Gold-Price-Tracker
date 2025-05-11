package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

func fetch() (float64,float64){
	resp,err := http.Get("https://data-asg.goldprice.org/dbXRates/USD")
  if err!= nil{
    fmt.Println("Error:",err)
    return 0,0
  }
  defer resp.Body.Close()

  body,err :=ioutil.ReadAll(resp.Body)
  if err!= nil{
    fmt.Println("Read Error:",err)
    return 0,0
  }

  var result struct{
      Items [] struct{
        XAUPrice float64 `json:"xauPrice"`
		XAGPrice float64 `json:"xagPrice"`
      }`json:"items"`
  }

  err = json.Unmarshal(body ,&result)
  if err!= nil{
    fmt.Println("JSON parse error:",err)
    return 0,0
  }

  fmt.Println("Price Gold:",result.Items[0].XAUPrice)
  fmt.Println("Price Silver:",result.Items[0].XAGPrice)
  return result.Items[0].XAUPrice , result.Items[0].XAGPrice
}

func main() {
	var chatIDs = make(map[int64] bool)


	err:= godotenv.Load()
	if err != nil{
		log.Fatal("Error Loading .env File")
	}

	botToken :=os.Getenv("API")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go func() {
		updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	
		for update := range updates {
			if update.Message != nil {
				id := update.Message.Chat.ID
				chatIDs[id] = true // track this chat
				log.Printf("New chat ID: %d", id)
			}
		}
	}()

	for{
				goldPrice,silverPrice := fetch()
				returnText := fmt.Sprintf("Gold Price Currently:%v \nSilver Price Currently:%v",goldPrice,silverPrice)

				for id:= range chatIDs{
				msg := tgbotapi.NewMessage(id,returnText)
				bot.Send(msg)	
				}
				time.Sleep(5 * time.Second)
		
	}
}

