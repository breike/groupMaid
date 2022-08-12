package main

import ( "fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidChatConfig(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	var msg_txt string = ""
	var err error      = nil

	args_list := strings.Split(update.Message.Text, " ")
	if len(args_list) > 1 {
		if args_list[1] == "get" {
			if len(args_list) > 2 {
				msg_txt, err = maidGetChatConfig(bot, update, db, args_list[2])
				if err != nil {
					msg_txt = "ERROR: can't get chat config, see logs for further info."

					return msg_txt, err
				}
			} else {
				msg_txt, err = maidGetChatConfig(bot, update, db, "")
				if err != nil {
					msg_txt = "ERROR: can't get chat config, see logs for further info."

					return msg_txt, err
				}
			}
		} else if args_list[1] == "set" {
			if len(args_list) > 3 {
				msg_txt, err = maidSetChatConfig(bot, update, db, args_list[2], args_list[3])
				if err != nil {
					msg_txt = "ERROR: can't set chat config, see logs for further info."

					return msg_txt, err
				}
			} else {
				msg_txt = "ERROR: /config set syntax: `/config set %key% %value%`"

				return msg_txt, err
			}
		}
	} else {
		msg_txt = "ERROR: need `set` or `get` command for config."

		return msg_txt, err
	}

	return msg_txt, err
}

func maidGetChatConfig(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB, key string) (string, error) {
	var msg_txt string = ""
	var err error      = nil

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
		msg_txt = "Config for: " + chat_name + " (" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")"
		msg_txt = msg_txt + "\nban\\_command\\: " + db.Chats[chat_id].Config.BanCmd
		msg_txt = msg_txt + "\nconfig\\_command\\: " + db.Chats[chat_id].Config.ConfigCmd
		msg_txt = msg_txt + "\nhelp\\_command\\: " + db.Chats[chat_id].Config.HelpCmd
		msg_txt = msg_txt + "\ninfo\\_command\\: " + db.Chats[chat_id].Config.InfoCmd
		msg_txt = msg_txt + "\nkick\\_command\\: " + db.Chats[chat_id].Config.KickCmd
		msg_txt = msg_txt + "\nmute\\_command\\: " + db.Chats[chat_id].Config.MuteCmd
		msg_txt = msg_txt + "\nremove\\_command\\: " + db.Chats[chat_id].Config.RemoveCmd
		msg_txt = msg_txt + "\nrules\\_command\\: " + db.Chats[chat_id].Config.RulesCmd
		msg_txt = msg_txt + "\nset\\_command\\: " + db.Chats[chat_id].Config.SetCmd
		msg_txt = msg_txt + "\nsetrules\\_command\\: " + db.Chats[chat_id].Config.SetrulesCmd
		msg_txt = msg_txt + "\nsetwelcome\\_command\\: " + db.Chats[chat_id].Config.SetwelcomeCmd
		msg_txt = msg_txt + "\nstatus\\_command\\: " + db.Chats[chat_id].Config.StatusCmd
		msg_txt = msg_txt + "\nunmute\\_command\\: " + db.Chats[chat_id].Config.UnmuteCmd
		msg_txt = msg_txt + "\nunset\\_command\\: " + db.Chats[chat_id].Config.UnsetCmd
		msg_txt = msg_txt + "\nupdate\\_command\\: " + db.Chats[chat_id].Config.UpdateCmd
		msg_txt = msg_txt + "\nwarn\\_command\\: " + db.Chats[chat_id].Config.WarnCmd
		msg_txt = msg_txt + "\nwelcome\\_command\\: " + db.Chats[chat_id].Config.WelcomeCmd
		msg_txt = msg_txt + "\nremove\\_command\\: " + db.Chats[chat_id].Config.RemoveCmd
		msg_txt = msg_txt + "\nban\\_command\\_on: " + strconv.FormatBool(db.Chats[chat_id].Config.BanCommandOn)
		msg_txt = msg_txt + "\ndelete\\_last\\_welcome: " + strconv.FormatBool(db.Chats[chat_id].Config.DeleteLastWelcome)
		msg_txt = msg_txt + "\ndisable\\_web\\_page\\_preview: " + strconv.FormatBool(db.Chats[chat_id].Config.DisableWebPagePreview)
		msg_txt = msg_txt + "\nwelcome_\\disable\\_web\\_page\\_preview: " + strconv.FormatBool(db.Chats[chat_id].Config.WelcomeDisableWebPagePreview)
		msg_txt = msg_txt + "\nrules_\\disable\\_web\\_page\\_preview: " + strconv.FormatBool(db.Chats[chat_id].Config.RulesDisableWebPagePreview)
		msg_txt = msg_txt + "\nhelp\\_command\\_on: " + strconv.FormatBool(db.Chats[chat_id].Config.HelpCommandOn)
		msg_txt = msg_txt + "\ninfo\\_command\\_on: " + strconv.FormatBool(db.Chats[chat_id].Config.InfoCommandOn)
		msg_txt = msg_txt + "\nmute\\_command\\_on: " + strconv.FormatBool(db.Chats[chat_id].Config.MuteCommandOn)
		msg_txt = msg_txt + "\nrules\\_command\\_on: " + strconv.FormatBool(db.Chats[chat_id].Config.RulesCommandOn)
		msg_txt = msg_txt + "\nwelcome\\_on: " + strconv.FormatBool(db.Chats[chat_id].Config.WelcomeOn)
		msg_txt = msg_txt + "\nwarns\\_limit: " + strconv.Itoa(db.Chats[chat_id].Config.WarnsLimit)
		msg_txt = msg_txt + "\nwarns\\_action: " + strconv.Itoa(db.Chats[chat_id].Config.WarnsAction) + " (" + warns_action + ")"
	}

	return msg_txt, err
}

func maidSetChatConfig(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB, key string, value string) (string, error) {
	var msg_txt string = ""
	var err error      = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)

	switch key {
	case "ban_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.BanCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }

		    msg_txt = "command has been changed"
		}
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "config_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.ConfigCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }

		    msg_txt = "command has been changed"
		}
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "help_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.HelpCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "info_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.InfoCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "kick_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.KickCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "mute_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.MuteCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "remove_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.RemoveCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "rules_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.RulesCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "set_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.SetCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "setrules_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.SetrulesCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "setwelcome_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.SetwelcomeCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "status_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.StatusCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "unmute_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.UnmuteCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "unset_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.UnsetCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "update_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.UpdateCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "warns_action":
		warns_action, err := strconv.Atoi(value)
		if  err != nil {
			msg_txt = "ERROR: value is not digit"

			return msg_txt, err
		}

		db.Chats[chat_id].Config.WarnsAction = warns_action

		err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "warn_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.WarnCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
	case "warns_limit":
		warns_limit, err := strconv.Atoi(value)
		if  err != nil {
			msg_txt = "ERROR: value is not digit"

			return msg_txt, err
		}

		db.Chats[chat_id].Config.WarnsLimit = warns_limit

		err = dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	case "welcome_command":
		if value[0:1] != "/" {
		    msg_txt = "command must start with '/'"

			return msg_txt, err
		} else {
		    db.Chats[chat_id].Config.WelcomeCmd = value

		    err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		    if err != nil {
			    msg_txt = "ERROR: can't write db, check out logs for further info"

			    return msg_txt, err
		    }
		}

		msg_txt = "command has been changed"
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
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

		err := dbWriteChatConfig(chat_id, db.Chats[chat_id].Config)
		if err != nil {
			msg_txt = "ERROR: can't write db, check out logs for further info"

			return msg_txt, err
		}

		msg_txt = "value has been written"
	}

	return msg_txt, err
}
