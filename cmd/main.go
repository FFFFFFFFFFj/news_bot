package main

import (
	//"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	
	"telegram-news-bot/internal/db"
	"telegram-news-bot/internal/bot"
	"telegram-news-bot/internal/config"
)

func main() {
    //	start .env file 
	config.LoadEnv()

	// Connecting to the base
	log.Println("Trying to connect to DB...")
	pool, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("DB conection error: %v", err)
	} 
	defer pool.Close()
	log.Println("DB connection established in main.go")
	
	//Telegram bot
	botToken := os.Getenv("BOT_TOKEN")
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	//Enable debugging (we will see all messages in the console)
	botAPI.Debug = false // debug OFF
	//log.Printf("Logged in as %s", botAPI.Self.UserName)

	//Settings up receiving updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 
	
	updates := botAPI.GetUpdatesChan(u)

	// ID super-admin
	superAdminID := os.Getenv("SUPER_ADMIN_ID")
	
	//main loop
	for update := range updates {
		bot.Router(update, botAPI, pool, superAdminID)
	}
}
