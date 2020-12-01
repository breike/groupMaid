package main

import (
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidGetRules(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	rules_text := db.Chats[chat_id].Config.RulesMsg
	rules_text = strings.Replace(rules_text, "\\```", "```", -1)
	rules_text = strings.Replace(rules_text, "\\`", "`", -1)
	rules_text = strings.Replace(rules_text, "$name", update.Message.From.FirstName, -1)

	msg_txt = rules_text

	return msg_txt, err
}

func maidSetRules(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error){
	var msg_txt string = ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	rules_text := strings.Replace(update.Message.Text, "/setrules ", "", -1)

	db.Chats[chat_id].Config.RulesMsg = rules_text

	err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
	if err != nil {
		return msg_txt, err
	}

	msg_txt = "Rules has been written"

	return msg_txt, err
}
