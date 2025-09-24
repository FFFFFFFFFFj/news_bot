package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"telegram-news-bot/internal/db"
)

//start bot
func HandleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool, superAdminID string) {
	role := "user"
	if fmt.Sprintf("%d", update.Message.From.ID) == superAdminID {
		role = "admin"
	}

	err := db.AddUserWithRoleIfNotExists(pool, update.Message.From.ID, update.Message.From.UserName, role)
	if err != nil {
		log.Printf("DB error: %v", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я бот парсер новостных статей в мире инвестиций." +
													   "\n\t/profile - открыть профиль" +
													   "\n\t/help - список команд бота")
	bot.Send(msg)
}

// command /help
func HandleHelp(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your commands: \n\t/start\n\thelp\n\t/addsource <name> <url>")
	bot.Send(msg)
}

// command for admins onli
func HandleAddSource(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool, userRole string) {
	if userRole != "admin" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This command is only available to administrators.")
		bot.Send(msg)
		return
	}

	userID := update.Message.From.ID
	state := GetUserState(userID)
	state.Current = StateAwaitingName

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter the source name: ")
	bot.Send(msg)
}

//message processing for status
func HandleState(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool, userRole string) {
	userID := update.Message.From.ID
	state := GetUserState(userID)

	switch state.Current {
	case StateAwaitingName:
		state.TempName = update.Message.Text
		state.Current = StateAwaitingURL
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "New enter the source URL: ")
		bot.Send(msg)
		
	case StateAwaitingURL:
		url := update.Message.Text
		name := state.TempName

		err := db.AddSource(pool, name, url)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error adding source:", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Source added: %s (%s)", name, url))
			bot.Send(msg)
		}

		ResetUserState(userID)
	}
}
