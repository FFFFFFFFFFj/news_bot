package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"telegram-news-bot/internal/db"
)

func main() {
	//start .env file 
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	//Enable debugging (we will see all messages in the console)
	bot.Debug = false // debug OFF
	log.Printf("Logged in as %s", bot.Self.UserName)

	//Settings up receiving updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 

	superAdminID := os.Getenv("SUPER_ADMIN_ID")
	
	updates := bot.GetUpdatesChan(u)

	//main loop
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("New message: %s", update.Message.Text)

		switch update.Message.Text {
		case "/start":
			role := "user"
			if fmt.Sprintf("%d", update.Message.From.ID) == superAdminID {
				role = "admin"
			}
			
			err := db.AddUserWithRoleIfNotExists(pool, update.Message.From.ID, update.Message.From.UserName, role)
			if err != nil {
				log.Printf("DB error: %v", err)
			} 
			
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, I`m your news bot")
			bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your command ....\n\t/start\n\t/help")
			bot.Send(msg)		
		}
	}
}
