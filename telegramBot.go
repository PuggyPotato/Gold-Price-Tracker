package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func parseInterval(input string) (int,error){
		input = strings.ToLower(strings.TrimSpace(input))
		if strings.HasSuffix(input,"min"){
			numStr := strings.TrimSuffix(input,"min")
			num,err :=strconv.Atoi(strings.TrimSpace(numStr))
			return num * 60,err
		}else if strings.HasSuffix(input,"hour"){
			numStr := strings.TrimSuffix(input,"hour")
			num,err :=strconv.Atoi(strings.TrimSpace(numStr))
			return num * 3600,err
		}else if strings.HasSuffix(input,"day"){
			numStr := strings.TrimSuffix(input,"day")
			num,err :=strconv.Atoi(strings.TrimSpace(numStr))
			return num * 86400,err
		}else{
			num,err := strconv.Atoi(input)
			return num,err
		}
}

func parseTarget(input string) (int,bool,bool,error){
	input = strings.ToLower(strings.TrimSpace(input))
		if strings.Contains(input,"silver"){
			if strings.Contains(input,"exceed"){
				numStr := strings.TrimSpace(strings.Replace(input, "silver exceed", "", 1))
				num,err :=strconv.Atoi(strings.TrimSpace(numStr))
				return num ,true,false ,err
			}else if strings.HasSuffix(input,"below"){
				numStr := strings.TrimSpace(strings.Replace(input, "silver below", "", 1))
				num,err :=strconv.Atoi(strings.TrimSpace(numStr))
				return num ,false,false,err
			}else{
			num,err := strconv.Atoi(input)
			return num,false,false,err
			}
		}else if strings.Contains(input,"gold"){
			if strings.Contains(input,"exceed"){
				numStr := strings.TrimSpace(strings.Replace(input, "gold exceed", "", 1))
				num,err :=strconv.Atoi(strings.TrimSpace(numStr))
				return num ,true,true ,err
			}else if strings.Contains(input,"below"){
				numStr := strings.TrimSpace(strings.Replace(input, "gold below", "", 1))
				num,err :=strconv.Atoi(strings.TrimSpace(numStr))
				return num ,false,true,err
			}else{
			num,err := strconv.Atoi(input)
			return num,false,true,err
			}
		}else{
			num,err := strconv.Atoi(input)
			return num,false,true,err
		}
		
}


func main() {
	var chatIDs = make(map[int64] bool)
	var waitingForInterval = make(map[int64]bool)
	var userIntervals = make(map[int64]int)
	var targetPrice = make(map[int64] int)
	var waitingForPrice = make(map[int64] bool)


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

				
				userInput := update.Message.Text
				userInput = strings.ToLower(userInput)

				if userInput == "/start"{
					reply := "Welcome To Potato Gold Bot, A Bot To Track Gold/Silver Prices, You can select how much message you get per interval,or set a price to track and notify you if the price exceeds or drop below the treshold. Start with \"gold\" or \"silver\" or \"price\" or /interval or /setTarget"
					replyMessage := tgbotapi.NewMessage(id,reply)
					bot.Send(replyMessage)
				}else if userInput == "gold" || userInput == "gold price" {
					goldPrice,_ := fetch()
					reply := fmt.Sprintf("Gold Price is currently:%.2f USD",goldPrice)
					replyMessage := tgbotapi.NewMessage(id,reply)
					bot.Send(replyMessage)
				}else if userInput == "silver" || userInput == "silver price" {
					_,silverPrice := fetch()
					reply := fmt.Sprintf("Silver Price is currently:%.2f USD",silverPrice)
					replyMessage := tgbotapi.NewMessage(id,reply)
					bot.Send(replyMessage)
				}else if userInput == "price"{
					goldPrice,silverPrice := fetch()
					reply := fmt.Sprintf("Gold Price is currently:%.2f USD \nSilver Price Is Currently %.2f USD",goldPrice,silverPrice)
					replyMessage := tgbotapi.NewMessage(id,reply)
					bot.Send(replyMessage)
				}else if userInput == "/interval"{
					replyMessage := tgbotapi.NewMessage(id,"What Is The Interval You Would Like To Receive The Update? Example:10min 2hour 1day")
					bot.Send(replyMessage)
					waitingForInterval[id] = true
				}else if waitingForInterval[id]{
					interval,err := parseInterval(userInput)
					if err !=nil{
						msg := tgbotapi.NewMessage(id,"Invalid,Try Again.")
						bot.Send(msg)
					}else{
						userIntervals[id] = interval
						waitingForInterval[id] = false

						msg := tgbotapi.NewMessage(id,"Succesfully Set Interval.")
						bot.Send(msg)
						

						go func(chatID int64,interval int){
								for{
								goldPrice,silverPrice := fetch()
								returnText := fmt.Sprintf("Gold Price Currently:%v \nSilver Price Currently:%v",goldPrice,silverPrice)

								for id:= range chatIDs{
								msg := tgbotapi.NewMessage(id,returnText)
								bot.Send(msg)	
								}
								time.Sleep(time.Duration(interval) * time.Second)
								}
						}(id,interval)
						
					}
				}else if userInput == "/settarget"{
					replyMessage := tgbotapi.NewMessage(id,"What Is The Targeted Price? eg. Silver exceed 3000...Gold below 3000 ")
					bot.Send(replyMessage)
					waitingForPrice[id] = true
				}else if waitingForPrice[id]{
					target,above,gold,err := parseTarget(userInput)
					if err != nil{
						msg := tgbotapi.NewMessage(id,"Invalid,Try Again.")
						bot.Send(msg)
					}else{
						targetPrice[id] = target
						waitingForPrice[id] = false

						msg := tgbotapi.NewMessage(id,"Succesfully Set Target.")
						bot.Send(msg)
						

						go func(chatID int64,above bool,gold bool,target int){
								for{
									goldPrice,silverPrice := fetch()
								if goldPrice > float64(target) && gold && above{
									returnText := fmt.Sprintf("Gold Price Exceeded %v \nGold Price Currently:%v",target,goldPrice)
									msg := tgbotapi.NewMessage(id,returnText)
									bot.Send(msg)	
								}else if goldPrice < float64(target) && gold && !above{
									returnText := fmt.Sprintf("Gold Price Dropped Below %v \nGold Price Currently:%v",target,goldPrice)
									msg := tgbotapi.NewMessage(id,returnText)
									bot.Send(msg)
								}
								if silverPrice > float64(target) && !gold && above{
									returnText := fmt.Sprintf("Silver Price Exceeded %v \nSilver Price Currently:%v",target,silverPrice)
									msg := tgbotapi.NewMessage(id,returnText)
									bot.Send(msg)	
								}else if silverPrice < float64(target) && !gold && !above{
									returnText := fmt.Sprintf("Silver Price Dropped Below %v \nSilver Price Currently:%v",target,silverPrice)
									msg := tgbotapi.NewMessage(id,returnText)
									bot.Send(msg)	
								}
								time.Sleep(60 * time.Second)

						}
						}(id,above,gold,target) 
						
					}
					
				}else if userInput == "/stop"{
					delete(chatIDs, id)
					delete(waitingForInterval, id)
					delete(userIntervals, id)
					delete(targetPrice, id)
					delete(waitingForPrice, id)
					msg := tgbotapi.NewMessage(id,"Succesfully Cleared.")
					bot.Send(msg)
				}else{
					replyMessage := tgbotapi.NewMessage(id,"I Dont Understand, Try \"gold\" or \"silver\" or \"price\"")
					bot.Send((replyMessage))
				}

				
				//reply := tgbotapi.NewMessage(id, "You said: " + userInput)
				
			}
		}
	}()



/*	for{
				goldPrice,silverPrice := fetch()
				returnText := fmt.Sprintf("Gold Price Currently:%v \nSilver Price Currently:%v",goldPrice,silverPrice)

				for id:= range chatIDs{
				msg := tgbotapi.NewMessage(id,returnText)
				bot.Send(msg)	
				}
				time.Sleep(60 * time.Second)
		
	}*/
	select{}
}