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

		if len(config.BotWhitelist) > 0 {
			var isWhitelisted = false

			for i := 0; i < len(config.BotWhitelist); i++ {
				if update.Message.Chat.ID == config.BotWhitelist[i] {
					isWhitelisted = true
				}
			}

			if !(isWhitelisted == true) {
				msg_text := "Whoops! Your chat is not in whitelist!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_text)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("ERROR: ", err)
				}

				chat_to_leave := tgbotapi.ChatConfig{ChatID: update.Message.Chat.ID}
				bot.LeaveChat(chat_to_leave)
				continue
			}
		}

		//check if no the chat settings in the db
		chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
		if _, ok := db.Chats[chat_id]; !(ok) {
			var chat_cfg = Chat_cfg_defaults

			db.Chats[chat_id] = new(chat)
			db.Chats[chat_id].Config = chat_cfg

			err := dbWriteChatConfig(chat_id, chat_cfg)
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
			resp, err := bot.Send(msg)
			if err != nil {
				log.Println("ERROR: ", err)

				continue
			}

			log.Println("LOG: message sent: ", resp)

			db.Chats[chat_id].Config.LastWelcomeID = resp.MessageID
			err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
			if err != nil {
				log.Fatal("ERROR: can't write chat config: ", err)
			}

			continue
		}

		/*
		memberFromCmd, err := maidGetChatMember(bot, update)
		if err != nil {
			msg.Text = "ERROR: internal error check log for the further info"
		}
		*/

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "ban":
				has_privileges, err := maidIsUserHasPrivileges(50, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidBanUser(bot, update, &db)
				}
			case "config":
				has_privileges, err := maidIsUserHasPrivileges(100, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidChatConfig(bot, update, &db)
					if err != nil {
						log.Println("ERROR: some problems with maidChatConfig: ", err)
					}
				}

			case "help":
				msg.Text, err = "type /hey", nil
			case "info":
				msg.Text, err = maidGetUserInfo(bot, update, &db)
				if err != nil {
					log.Println("ERROR: Failed to unset user info: ", err)
				}
			case "kick":
				has_privileges, err := maidIsUserHasPrivileges(50, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidKickUser(bot, update)
				}
				if err != nil {
					log.Println("ERROR: Failed to kick user: ", err)
				}
			case "mute":
				has_privileges, err := maidIsUserHasPrivileges(50, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidMuteUser(bot, update)
				}
			case "remove":
				has_privileges, err := maidIsUserHasPrivileges(50, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidRemoveUserInfo(bot, update, &db)
				}
				if err != nil {
					log.Println("ERROR: Failed to remove user data: ", err)
				}
			case "rules":
				msg.Text, err = maidGetRules(bot, update, &db)
				msg.DisableWebPagePreview = db.Chats[chat_id].Config.RulesDisableWebPagePreview

			case "set":
				has_privileges, err := maidIsUserHasPrivileges(10, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidSetUserInfo(bot, update, &db)
					if err != nil {
						log.Println("ERROR: Failed to set user info: ", err)
					}
				}
			case "setrules":
				has_privileges, err := maidIsUserHasPrivileges(70, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidSetRules(bot, update, &db)
				}
			case "setwelcome":
				has_privileges, err := maidIsUserHasPrivileges(70, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidSetWelcome(bot, update, &db)
					if err != nil {
						log.Println("ERROR: Failed to set welcome message: ", err)
					}
				}
			case "status":
				has_privileges, err := maidIsUserHasPrivileges(50, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidGetUserStatus(bot, update, &db)
				}
				if err != nil {
					log.Println("ERROR: Failed to reset warns: ", err)
				}
			case "unmute":
				has_privileges, err := maidIsUserHasPrivileges(100, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidUnmuteUser(bot, update)
				}
			case "unset":
				has_privileges, err := maidIsUserHasPrivileges(10, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidUnsetUserInfo(bot, update, &db)
					if err != nil {
						log.Println("ERROR: Failed to unset user info: ", err)
					}
				}
			case "update":
				has_privileges, err := maidIsUserHasPrivileges(100, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidUpdate(bot, update, &db)
					if err != nil {
						log.Println("ERROR: failed to warn user: ", err)
					}
				}
			case "warn":
				has_privileges, err := maidIsUserHasPrivileges(50, bot, update, &db)
				if err != nil {
					log.Println("ERROR: can't check user privileges: ", err)
				}

				if !(has_privileges) {
					msg.Text = "ERROR: you do not have needed privileges"
				} else {
					msg.Text, err = maidWarnUser(bot, update, &db)
					if err != nil {
						log.Println("ERROR: failed to warn user: ", err)
					}
				}
			case "welcome":
				if config.BotDebug {
					msg.Text, err = maidGetWelcome(bot, update, &db)
					msg.DisableWebPagePreview = db.Chats[chat_id].Config.WelcomeDisableWebPagePreview
					resp, err := bot.Send(msg)
					if err != nil {
						log.Println("ERROR: ", err)
					}

					log.Println("LOG: message sent: ", resp)

					db.Chats[chat_id].Config.LastWelcomeID = resp.MessageID
					err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
					if err != nil {
						log.Fatal("ERROR: can't write chat config: ", err)
					}

					continue
				}
			}

			if err != nil {
				log.Println("ERROR: ", err)
			}

			if msg.Text != "" {
				resp, err := bot.Send(msg)
				if err != nil {
					log.Println("ERROR: ", err)
				} else {
					log.Println("LOG: message sent: ", resp)
				}
			}
		}
	}
}
