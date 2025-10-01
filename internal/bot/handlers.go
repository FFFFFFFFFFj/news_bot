package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"telegram-news-bot/internal/db"
	"telegram-news-bot/internal/texts"
)

//start bot
func HandleStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool, superAdminID string) {
	role := "user"
	
	//If this is a super-admin, we set the role to superadmin
	if fmt.Sprintf("%d", update.Message.From.ID) == superAdminID {
		role = "superadmin"
	}
	
	//Add the user to the database if it doesn't exist
	err := db.AddUserWithRoleIfNotExists(pool, update.Message.From.ID, update.Message.From.UserName, role)
	if err != nil {
		log.Printf("DB error: %v", err)
	}
	
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.StartMessage)
	bot.Send(msg)
}

// command /help
func HandleHelp(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.HelpMessage)
	bot.Send(msg)
}

// command for admins onli
func HandleAddSource(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool, userRole string) {
	if userRole != "admin" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, texts.AdminOnly)
		bot.Send(msg)
		return
	}

	userID := update.Message.From.ID
	state := GetUserState(userID)
	state.Current = StateAwaitingName

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название источника: ")
	bot.Send(msg)
}

// comand /linkchannel
func HandleLinkChannel(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	state := GetUserState(userID)
	state.Current = StateAwaitingChannelName

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название идентификатор канала @name")
	bot.Send(msg)
}

// comand /unlinkchannel
func HandleUnlinkChannel(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	state := GetUserState(userID)
	state.Current = StateAwaitingUnlinkChannel

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название идентификатора канала @name")
	bot.Send(msg)
}

// comand /setposttime
func HandleSetPostTime(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	state := GetUserState(userID)
	state.Current = StateAwaitingPostTime

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите время запланированного постинга в формате: (09:00)")
	bot.Send(msg)
}

// comand /setpostcount
func HandleSetPostCount(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userID := update.Message.From.ID
	state := GetUserState(userID)
	state.Current = StateAwaitingPostCount

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите колличество постов в день:")
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пока не знаю что это делает, но тут нужен url: ")
		bot.Send(msg)
		
	case StateAwaitingURL:
		url := update.Message.Text
		name := state.TempName

		err := db.AddSource(pool, name, url)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка добавления источника:", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Источник добавлен: %s (%s)", name, url))
			bot.Send(msg)
		}
		ResetUserState(userID)

	case StateAwaitingChannelName:
		channel := update.Message.Text
		if err := db.AddUserChannel(pool, userID, channel); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Канал  %s успешно привязан!", channel))
			bot.Send(msg)
		}
		ResetUserState(userID)

	case StateAwaitingUnlinkChannel:
		channel := update.Message.Text
		if err := db.RemoveUserChannel(pool, userID, channel); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Канал  %s оключен.", channel))
			bot.Send(msg)
		}
		ResetUserState(userID)

	case StateAwaitingPostTime:
		time := update.Message.Text
		if err := db.UpdateUserChannelTime(pool, userID, time); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Время публикации установлено: %s", time))
			bot.Send(msg)
		}
		ResetUserState(userID)

	case StateAwaitingPostCount:
		count := update.Message.Text
		var c int
		_, err := fmt.Sscanf(count, "%d", &c)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста введите коррктный номер.")
			bot.Send(msg)
			return
		}
		if err := db.UpdateUserChannelCount(pool, userID, c); err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Количество отправленых постов: %d"))
			bot.Send(msg)
		}
		ResetUserState(userID)
	}
}

//Profile
func HandleProfile(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool) {
	var role, username string
	var telegramID int64

	err := pool.QueryRow(
		context.Background(),
		"SELECT telegram_id, username, role FROM users WHERE telegram_id=$1",
		update.Message.From.ID,
	).Scan(&telegramID, &username, &role)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка получения профиля: %v", err))
		bot.Send(msg)
		return
	}
	//If username is not empty we put a placeholder
	if username == "" {
		username = "not specified"
	}

	msgText := fmt.Sprintf(texts.ProfileTemplate, telegramID, username, role)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}

func HandleListSources(update tgbotapi.Update, bot *tgbotapi.BotAPI, pool *pgxpool.Pool) {
	rows, err := pool.Query(context.Background(), "SELECT name, url FROM sources ORDER BY id")
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка получения источников: %v", err))
		bot.Send(msg)
		return
	}
	defer rows.Close()

	msgText := "Список источников:\n"
	for rows.Next() {
		var name, url string
		if err := rows.Scan(&name, &url); err != nil {
			continue
		}
		msgText += fmt.Sprintf("- %s (%s)\n", name, url)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	bot.Send(msg)
}
