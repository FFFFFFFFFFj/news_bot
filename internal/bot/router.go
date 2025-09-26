package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Router(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool, superAdminID string) {
	if update.Message == nil {
		return
	}

	//Obtaining a user role (user/admin)
	userRole := "user"
	if fmt.Sprintf("%d", update.Message.From.ID) == superAdminID {
		userRole = "admin"
	}

	// status check
	state := GetUserState(update.Message.From.ID)
	if state.Current != StateNone && update.Message.Text[0] != '/' {
		HandleState(update, bot, pool, userRole)
		return
	}

	switch update.Message.Text {
	case "/start":
		HandleStart(update, bot, pool, superAdminID)
	case "/help":
		HandleHelp(update, bot)
	case "/addsource":
		HandleAddSource(update, bot, pool, userRole)
	case "/linkchannel":
		HandleLinkChannel(update, bot)
	case "/unlinkchannel":
		HandleUnlinkChannel(update, bot)
	case "/setposttime":
		HandleSetPostTime(update, bot)
	case "/setpostcount":
		HandleSetPostCount(update, bot)
	}
}
