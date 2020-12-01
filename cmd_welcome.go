package main

import (
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidGetWelcome(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	welcome_text := db.Chats[chat_id].Config.WelcomeMsg
	welcome_text = strings.Replace(welcome_text, "\\```", "```", -1)
	welcome_text = strings.Replace(welcome_text, "\\`", "`", -1)
	welcome_text = strings.Replace(welcome_text, "$name", update.Message.From.FirstName, -1)

	msg_txt = welcome_text

	return msg_txt, err
}

func maidSetWelcome(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	welcome_text := strings.Replace(update.Message.Text, "/setwelcome ", "", -1)

	db.Chats[chat_id].Config.WelcomeMsg = welcome_text

	err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
	if err != nil {
		return msg_txt, err
	}

	msg_txt = "Welcome has been written"

	return msg_txt, err

}
