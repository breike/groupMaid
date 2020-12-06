package main

import (
	"fmt"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidGetChatConfig(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB, key string) (string, error) {
	msg_txt         := ""
	var err error    = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)

	var warns_action string;
	switch db.Chats[chat_id].Config.WarnsAction {
	case 0:
		warns_action = "ban"
	case 1:
		warns_action = "kick"
	case 2:
		warns_action = "mute"
	}

	var chat_name string;
	if update.Message.Chat.Title != "" {
		chat_name = update.Message.Chat.Title
	} else {
		chat_name = "chat"
	}

	if key != "" {
		switch key {
		case "ban_command_on":
			msg_txt = fmt.Sprintf("ban\\_command\\_on: %t",
								  db.Chats[chat_id].Config.BanCommandOn)
		case "delete_last_welcome":
			msg_txt = fmt.Sprintf("delete\\_last\\_welcome: %t",
								  db.Chats[chat_id].Config.DeleteLastWelcome)
		case "disable_web_page_preview":
			msg_txt = fmt.Sprintf("disable\\_web\\_page\\_preview: %t",
								  db.Chats[chat_id].Config.DisableWebPagePreview)
		case "welcome_disable_web_page_preview":
			msg_txt = fmt.Sprintf("welcome\\_disable\\_web\\_page\\_preview: %t",
								  db.Chats[chat_id].Config.WelcomeDisableWebPagePreview)
		case "rules_disable_web_page_preview":
			msg_txt = fmt.Sprintf("rules\\_disable\\_web\\_page\\_preview: %t",
								  db.Chats[chat_id].Config.RulesDisableWebPagePreview)
		case "help_command_on":
			msg_txt = fmt.Sprintf("help\\_command\\_on: %t",
								  db.Chats[chat_id].Config.HelpCommandOn)
		case "info_command_on":
			msg_txt = fmt.Sprintf("info\\_command\\_on: %t",
								  db.Chats[chat_id].Config.InfoCommandOn)
		case "mute_command_on":
			msg_txt = fmt.Sprintf("mute\\_command\\_on: %t",
								  db.Chats[chat_id].Config.MuteCommandOn)
		case "rules_command_on":
			msg_txt = fmt.Sprintf("rules\\_command\\_on: %t",
								  db.Chats[chat_id].Config.RulesCommandOn)
		case "welcome_on":
			msg_txt = fmt.Sprintf("welcome\\_on: %t",
								  db.Chats[chat_id].Config.WelcomeOn)
		case "warns_limit":
			msg_txt = fmt.Sprintf("warns\\_limit: %d",
								  db.Chats[chat_id].Config.WarnsLimit)
		case "warns_action":
			msg_txt = fmt.Sprintf("warns\\_action: %d (%s)",
								  db.Chats[chat_id].Config.WarnsLimit,
								  warns_action)
		}
	} else {
		msg_txt = fmt.Sprintf("Config for %s:\n" +
							  "ban\\_command\\_on: %t\n" +
							  "delete\\_last\\_welcome: %t\n" +
							  "disable\\_web\\_page\\_preview: %t\n" +
							  "welcome\\_disable\\_web\\_page\\_preview: %t\n" +
							  "rules\\_disable\\_web\\_page\\_preview: %t\n" +
							  "help\\_command\\_on: %t\n" +
							  "info\\_command\\_on: %t\n" +
							  "mute\\_command\\_on: %t\n" +
							  "rules\\_command\\_on: %t\n" +
							  "welcome\\_on: %t\n" +
							  "warns\\_limit: %d\n" +
							  "warns\\_action: %d (%s)\n",
							  chat_name,
							  db.Chats[chat_id].Config.BanCommandOn,
							  db.Chats[chat_id].Config.DeleteLastWelcome,
							  db.Chats[chat_id].Config.DisableWebPagePreview,
							  db.Chats[chat_id].Config.WelcomeDisableWebPagePreview,
							  db.Chats[chat_id].Config.RulesDisableWebPagePreview,
							  db.Chats[chat_id].Config.HelpCommandOn,
							  db.Chats[chat_id].Config.InfoCommandOn,
							  db.Chats[chat_id].Config.MuteCommandOn,
							  db.Chats[chat_id].Config.RulesCommandOn,
							  db.Chats[chat_id].Config.WelcomeOn,
							  db.Chats[chat_id].Config.WarnsLimit,
							  db.Chats[chat_id].Config.WarnsAction, warns_action)
	}

	return msg_txt, err
}

func maidSetChatConfig(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB, key string, value string) (string, error) {
	msg_txt         := ""
	var err error    = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)

	switch key {
	case "ban_command_on":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.BanCommandOn = true
		} else {
			db.Chats[chat_id].Config.BanCommandOn = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "delete_last_welcome":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.DeleteLastWelcome = true
		} else {
			db.Chats[chat_id].Config.DeleteLastWelcome = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "disable_web_page_preview":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.DisableWebPagePreview = true
		} else {
			db.Chats[chat_id].Config.DisableWebPagePreview = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "welcome_disable_web_page_preview":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.WelcomeDisableWebPagePreview = true
		} else {
			db.Chats[chat_id].Config.WelcomeDisableWebPagePreview = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "rules_disable_web_page_preview":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.RulesDisableWebPagePreview = true
		} else {
			db.Chats[chat_id].Config.RulesDisableWebPagePreview = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "help_command_on":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.HelpCommandOn = true
		} else {
			db.Chats[chat_id].Config.HelpCommandOn = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "info_command_on":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.InfoCommandOn = true
		} else {
			db.Chats[chat_id].Config.InfoCommandOn = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "mute_command_on":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.MuteCommandOn = true
		} else {
			db.Chats[chat_id].Config.MuteCommandOn = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "rules_command_on":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.RulesCommandOn = true
		} else {
			db.Chats[chat_id].Config.RulesCommandOn = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "welcome_on":
		if value != "true" && value != "false" {
			msg_txt = "ERROR: unknown value"

			return msg_txt, err
		}

		if value == "true" {
			db.Chats[chat_id].Config.WelcomeOn = true
		} else {
			db.Chats[chat_id].Config.WelcomeOn = false
		}

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "warns_limit":
		warns_limit, err := strconv.Atoi(value)
		if  err != nil {
			msg_txt = "ERROR: value is not digit"

			return msg_txt, err
		}

		db.Chats[chat_id].Config.WarnsLimit = warns_limit

		err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "warns_action":
		warns_action, err := strconv.Atoi(value)
		if  err != nil {
			msg_txt = "ERROR: value is not digit"

			return msg_txt, err
		}

		db.Chats[chat_id].Config.WarnsAction = warns_action

		err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config, db)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	}

	return msg_txt, err
}
