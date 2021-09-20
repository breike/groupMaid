package main

import (
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

	member, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: internal error check log for the further info"

		return msg_txt, err
	}

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(member.User.ID)

	if db.Chats[chat_id].Users[user_id] == nil {
		db.Chats[chat_id].Users[user_id] = new(user)
	}

	// user info variables definition
	preferred_name := "unknown"
	if db.Chats[chat_id].Users[user_id].PreferredName != "" {
		preferred_name = db.Chats[chat_id].Users[user_id].PreferredName
	}

	gender  := "unknown"
	if db.Chats[chat_id].Users[user_id].Gender != "" {
		gender  = db.Chats[chat_id].Users[user_id].Gender
	}

	pronouns := "unknown"
	if db.Chats[chat_id].Users[user_id].Pronouns != "" {
		pronouns = db.Chats[chat_id].Users[user_id].Pronouns
	}

	notes   := "unknown"
	if db.Chats[chat_id].Users[user_id].Notes != "" {
		notes = db.Chats[chat_id].Users[user_id].Notes
	}

	// msg_txt generation
	msg_txt = "```" + "\nИмя:          " + preferred_name +
	          "\nГендер:       " + gender +
			  "\nМестоимения:  " + pronouns

	if pronouns != "unknown" {
		msg_txt = msg_txt + "\nЗаметки:      " + notes
	}

	msg_txt = msg_txt + "\n```"

	return msg_txt, err
}

func maidGetUserStatus(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: Reply to user you want to get status"

		return msg_txt, err
	}

	member, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: internal error check log for the further info"

		return msg_txt, err
	}

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(member.User.ID)

	if db.Chats[chat_id].Users[user_id] == nil {
		db.Chats[chat_id].Users[user_id] = new(user)
	}

	// user status variables definition
	admin_notes   := "unknown"
	if db.Chats[chat_id].Users[user_id].AdminNotes != "" {
		admin_notes = db.Chats[chat_id].Users[user_id].AdminNotes
	}

	is_member       := "Нет"
	if member.IsMember() {
		is_member = "Да"
	}

	is_kicked       := "Нет"
	if member.WasKicked() {
		is_kicked = "Да"
	}

	privileges      := db.Chats[chat_id].Users[user_id].Privileges
	privileges_text := "member"
	switch privileges {
	case 0:
		privileges_text = "member"
	case 10:
		privileges_text = "can change info"
	case 50:
		privileges_text = "moderator"
	case 70:
		privileges_text = "can set rules and welcome"
	case 100:
		privileges_text = "full admin rights"
	default:
		privileges_text = "member"
	}

	ban_note        := "unknown"
	if db.Chats[chat_id].Users[user_id].BanNote != "" {
		ban_note = db.Chats[chat_id].Users[user_id].BanNote
	}

	ban_from        := "unknown"
	if db.Chats[chat_id].Users[user_id].BanFrom != "" {
		ban_from = db.Chats[chat_id].Users[user_id].BanFrom
	}

	warn_number     := db.Chats[chat_id].Users[user_id].Warns
	warn_limit      := db.Chats[chat_id].Config.WarnsLimit

	// msg_txt generation
	msg_txt = "```"
	msg_txt = msg_txt + "\nИмя:            " + member.User.FirstName
	msg_txt = msg_txt + "\nВ чате:         " + is_member
	msg_txt = msg_txt + "\nКикнут:         " + is_kicked
	msg_txt = msg_txt + "\nПривилегии:     " +
	                    strconv.Itoa(privileges) + " (" +
						privileges_text + ")"

	if ban_note != "unknown" {
		msg_txt = msg_txt + "\nБан:            " + ban_note
	}

	if ban_from != "unknown" {
		msg_txt = msg_txt + "\nБан от:         " + ban_from
	}

	if admin_notes != "unknown" {
		msg_txt = msg_txt + "\nЗаметки:        " + admin_notes
	}

	msg_txt = msg_txt + "\nПредупреждения: " +
	                    strconv.Itoa(warn_number) + "/" +
						strconv.Itoa(warn_limit)
	msg_txt = msg_txt + "\n```"

	return msg_txt, err
}

func maidSetUserInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: Reply to user you want to set info."

		return msg_txt, err
	}

	chat_id  := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id  := strconv.Itoa(update.Message.ReplyToMessage.From.ID)
	admin_id := strconv.Itoa(update.Message.From.ID)

	if db.Chats[chat_id].Users[user_id] == nil {
		db.Chats[chat_id].Users[user_id] = new(user)
	}

    if config.BotAdminID == admin_id {
        db.Chats[chat_id].Users[admin_id].Privileges = 100
    }

	key   := strings.Split(update.Message.Text, " ")[1]
	value := strings.Split(update.Message.Text, " " + key + " ")[1]

	switch key {
	case "adminnotes":
		if db.Chats[chat_id].Users[admin_id].Privileges < 50 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].AdminNotes    = value
	case "ban":
		if db.Chats[chat_id].Users[admin_id].Privileges < 70 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].BanNote        = value
	case "banfrom":
		if db.Chats[chat_id].Users[admin_id].Privileges < 70 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].BanFrom        = value
	case "gender":
		db.Chats[chat_id].Users[user_id].Gender         = value
	case "pronouns":
		db.Chats[chat_id].Users[user_id].Pronouns       = value
	case "privileges":
		if db.Chats[chat_id].Users[admin_id].Privileges < 100 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].Privileges, err = strconv.Atoi(value)
	case "name":
		db.Chats[chat_id].Users[user_id].PreferredName   = value
	case "notes":
		db.Chats[chat_id].Users[user_id].Notes           = value
	default:
		msg_txt = "Bad key for setting user info."

		return msg_txt, err
	}

	msg_txt = "User's info has been set."

	err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
	if err != nil {
		msg_txt = "ERROR: internal error, check log for further info."
		return msg_txt, err
	}

	return msg_txt, err
}

func maidUnsetUserInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(update.Message.ReplyToMessage.From.ID)
	admin_id := strconv.Itoa(update.Message.From.ID)

	key   := strings.Split(update.Message.Text, " ")[1]

    if db.Chats[chat_id].Users[admin_id] == nil {
        db.Chats[chat_id].Users[admin_id] = new(user)
    }

    if config.BotAdminID == admin_id {
        db.Chats[chat_id].Users[admin_id].Privileges = 100
    }

	switch key {
	case "adminnotes":
		if db.Chats[chat_id].Users[admin_id].Privileges < 70 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].AdminNotes    = ""
	case "ban":
		if db.Chats[chat_id].Users[admin_id].Privileges < 70 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].BanNote       = ""
	case "banfrom":
		if db.Chats[chat_id].Users[admin_id].Privileges < 70 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].BanFrom       = ""
	case "gender":
		db.Chats[chat_id].Users[user_id].Gender        = ""
	case "pronouns":
		db.Chats[chat_id].Users[user_id].Pronouns      = ""
	case "name":
		db.Chats[chat_id].Users[user_id].PreferredName = ""
	case "notes":
		db.Chats[chat_id].Users[user_id].Notes         = ""
	case "warns":
		if db.Chats[chat_id].Users[admin_id].Privileges < 70 {
			msg_txt = "ERROR: you do not have needed privileges for that."

			return msg_txt, err
		}
		db.Chats[chat_id].Users[user_id].Warns         = 0
	default:
		msg_txt = "Bad key for removing user info."

		return msg_txt, err
	}

	msg_txt = "User's info has been unset."

	err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
	if err != nil {
		msg_txt = "ERROR: internal error, check log for further info."
		return msg_txt, err
	}

	return msg_txt, err
}

// decrease some user info by 1 (if this info countable)
func maidRemoveUserInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(update.Message.ReplyToMessage.From.ID)

	key   := strings.Split(update.Message.Text, " ")[1]

	switch key {
	case "warn":
		if db.Chats[chat_id].Users[user_id].Warns == 0 {
			msg_txt = "User does not have any warn."
		} else {
			db.Chats[chat_id].Users[user_id].Warns -= 1
			msg_txt = "User's warn has been removed."

			err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
			if err != nil {
				msg_txt = "ERROR: can't write warn to db, check out logs."

				return msg_txt, err
			}
		}
	default:
		msg_txt = "Bad key for removing user info."
		return msg_txt, err
	}


	err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
	if err != nil {
		msg_txt = "ERROR: internal error, check log for further info"
		return msg_txt, err
	}

	return msg_txt, err
}
