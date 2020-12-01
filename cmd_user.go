package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidGetUserInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: Reply to user you want to get info"

		return msg_txt, err
	}

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(update.Message.ReplyToMessage.From.ID)

	if db.Chats[chat_id].Users[user_id] == nil {
		db.Chats[chat_id].Users[user_id] = new(user)
	}

	gender  := "Unknown"
	if db.Chats[chat_id].Users[user_id].Gender != "" {
		gender  = db.Chats[chat_id].Users[user_id].Gender
	}

	notes   := "Unknown"
	if db.Chats[chat_id].Users[user_id].Notes != "" {
		notes = db.Chats[chat_id].Users[user_id].Notes
	}

	preferred_name := "Unknown"
	if db.Chats[chat_id].Users[user_id].PreferredName != "" {
		preferred_name = db.Chats[chat_id].Users[user_id].PreferredName
	}

	ban_note := "Unknown"
	if db.Chats[chat_id].Users[user_id].BanNote == "" {
		msg_txt = fmt.Sprintf("Имя: %s\nГендер: %s\nЗаметки: %s", preferred_name,
		                      gender, notes)
	} else {
		ban_note = db.Chats[chat_id].Users[user_id].BanNote
		msg_txt = fmt.Sprintf("Имя: %s\nГендер: %s\nЗаметки: %s\nБан: %s",
		                      preferred_name, gender, notes, ban_note)
	}

	return msg_txt, err
}

func maidSetUserInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(update.Message.ReplyToMessage.From.ID)

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: Reply to user you want to set info"

		return msg_txt, err
	}

	if db.Chats[chat_id].Users[user_id] == nil {
		db.Chats[chat_id].Users[user_id] = new(user)
	}

	key   := strings.Split(update.Message.Text, " ")[1]
	value := strings.Split(update.Message.Text, " ")[2]

	switch key {
	case "ban":
		db.Chats[chat_id].Users[user_id].BanNote       = value
	case "gender":
		db.Chats[chat_id].Users[user_id].Gender        = value
	case "name":
		db.Chats[chat_id].Users[user_id].PreferredName = value
	case "notes":
		db.Chats[chat_id].Users[user_id].Notes         = value
	}

	msg_txt = "User's info has been setted"

	err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
	if err != nil {
		msg_txt = "ERROR: internal error, check log for further info"
		return msg_txt, err
	}

	return msg_txt, err
}

func maidUnsetUserInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(update.Message.ReplyToMessage.From.ID)

	key   := strings.Split(update.Message.Text, " ")[1]

	switch key {
	case "ban":
		db.Chats[chat_id].Users[user_id].BanNote       = ""
	case "gender":
		db.Chats[chat_id].Users[user_id].Gender        = ""
	case "name":
		db.Chats[chat_id].Users[user_id].PreferredName = ""
	case "notes":
		db.Chats[chat_id].Users[user_id].Notes         = ""
	}

	msg_txt = "User's info has been unsetted"

	return msg_txt, err
}
