package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	telegramBotToken string
)

func init() {

	flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
	flag.Parse()

	if telegramBotToken == "" {
		log.Print("-telegrambottoken is required")
		os.Exit(1)
	}
}

func main() {

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err)
	}

	for update := range updates {

		reply := "Не знаю что сказать"
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = `I can tell you the latest news about the game valorant or Tarkov, for this tell me: /Valorant or /tarkov`
		case "valorant":
			reply = ParseValorant()

		case "tarkov":

			reply = Tarkov() //86
		case "Tarkov":
			reply = TarkovParser()

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		bot.Send(msg)
	}
}
