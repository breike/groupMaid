package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var config maidConfig

func main() {
	configPath := flag.String("config", "config.toml", "config file path")
	flag.Parse()

	config, err := configInit(*configPath)
	if err != nil {
		log.Fatal("ERROR: cat't init bot config: ", err)
	}

	if config.BotDebug {
		fmt.Println("\nbot config:")
		fmt.Println("\tBotDebug: ", config.BotDebug)
		fmt.Println("\tTgBotAPI: ", config.TgBotAPI)
		fmt.Println("\tBotDirectory: ", config.BotDirectory, "\n")
	}

	bot, err := tgbotapi.NewBotAPI(config.TgBotAPI)
	if err != nil {
		log.Fatal("ERROR: can't init bot API: ", err)
	}

	db, err := dbInit()
	if err != nil {
		log.Fatal("ERROR: can't init bot database: ", err)
	}

	if config.BotDebug {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	log.Printf("%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("ERROR: can't get updates channel: ", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//check if no the chat settings in the db
		chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
		if _, ok := db.Chats[chat_id]; !(ok) {
			var chat_cfg = Chat_cfg_defaults

			err := dbWriteChatConfig(chat_id, chat_cfg, &db)
			if err != nil {
				log.Fatal("ERROR: can't write chat config: ", err)
			}
		}

		if db.Chats[chat_id].Users == nil {
			db.Chats[chat_id].Users = make(map[string]*user)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = db.Chats[chat_id].Config.DisableWebPagePreview

		if update.Message.NewChatMembers != nil {
			msg.Text, err = maidGetWelcome(bot, update, &db)
			msg.DisableWebPagePreview = db.Chats[chat_id].Config.WelcomeDisableWebPagePreview
		}

		memberFromCmd, err := maidGetChatMember(bot, update)
		if err != nil {
			msg.Text = "ERROR: internal error check log for the further info"
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "ban":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidBanUser(bot, update, &db)
				}
			case "help":
				msg.Text, err = "type /hey", nil
			case "info":
				msg.Text, err = maidGetUserInfo(bot, update, &db)
				if err != nil {
					log.Println("ERROR: Failed to unset user info: ", err)
				}
			case "kick":
				msg.Text, err = maidKickUser(bot, update)
				if err != nil {
					log.Println("ERROR: Failed to kick user: ", err)
				}
			case "mute":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidMuteUser(bot, update)
				}
			case "rules":
				msg.Text, err = maidGetRules(bot, update, &db)
				msg.DisableWebPagePreview = db.Chats[chat_id].Config.RulesDisableWebPagePreview

			case "set":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidSetUserInfo(bot, update, &db)
					if err != nil {
						log.Println("ERROR: Failed to set user info: ", err)
					}
				}

			case "setrules":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidSetRules(bot, update, &db)
				}
			case "setwelcome":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidSetWelcome(bot, update, &db)
					if err != nil {
						log.Println("ERROR: Failed to set welcome message: ", err)
					}
				}
			case "welcome":
				if config.BotDebug {
					msg.Text, err = maidGetWelcome(bot, update, &db)
					msg.DisableWebPagePreview = db.Chats[chat_id].Config.WelcomeDisableWebPagePreview
				}
			case "warn":
				msg.Text, err = maidWarnUser(bot, update, &db)
				if err != nil {
					log.Println("ERROR: failed to warn user: ", err)
				}
			case "unmute":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidUnmuteUser(bot, update)
				}
			case "unset":
				if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
					msg.Text = "ERROR: not admin"
				} else {
					msg.Text, err = maidUnsetUserInfo(bot, update, &db)
					if err != nil {
						log.Println("ERROR: Failed to unset user info: ", err)
					}
				}
			}

			if err != nil {
				log.Println("ERROR: ", err)
			}

			resp, err := bot.Send(msg)
			if err != nil {
				log.Println("ERROR: ", err)
			}

			log.Println("LOG: message sent: ", resp)
		}
	}
}
