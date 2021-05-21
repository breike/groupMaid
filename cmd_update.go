package main

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	admin_id := update.Message.From.ID

	key   := strings.Split(update.Message.Text, " ")[1]

	switch key {
	case "db":
		if admin_id != config.BotAdminID {
			msg_txt = "ERROR: you do not have needed privileges for that"
		}

		*db, err = dbInit()
		if err != nil {
			msg_txt = "ERROR: can't load db, see logs for further information"
			return msg_txt, err
		}

		msg_txt = "DB has been loaded."
	default:
		msg_txt = "ERROR: bad key for command"
	}

	return msg_txt, err
}
