package main

import (
	"fmt"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidWarnUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil
	
	memberToWarn, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: can't get chat member from reply, check out logs"

		return msg_txt, err
	}

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(memberToWarn.User.ID)

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: Reply to user you want to warn"

		return msg_txt, err
	}

	if db.Chats[chat_id].Users[user_id] == nil {
		db.Chats[chat_id].Users[user_id] = new(user)
	}

	db.Chats[chat_id].Users[user_id].Warns += 1

	if db.Chats[chat_id].Users[user_id].Warns == db.Chats[chat_id].Config.WarnsLimit {
		db.Chats[chat_id].Users[user_id].Warns = 0

		err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
		if err != nil {
			msg_txt = "ERROR: can't write warn to db, check out logs"

			return msg_txt, err
		}

		switch db.Chats[chat_id].Config.WarnsAction {
		case 0:
			_, err := maidBanUser(bot, update, db)
			if err != nil {
				msg_txt = "ERROR: can't ban user"

				return msg_txt, err
			}

			msg_txt = fmt.Sprintf("%s reached warn limit and has been banned",
		                          memberToWarn.User.FirstName)
		case 1:
			_, err := maidKickUser(bot, update)
			if err != nil {
				msg_txt = "ERROR: can't kick user"

				return msg_txt, err
			}

			msg_txt = fmt.Sprintf("%s reached warn limit and has been kicked",
		                          memberToWarn.User.FirstName)
		case 2:
			_, err := maidMuteUser(bot, update)
			if err != nil {
				msg_txt = "ERROR: can't mute user"

				return msg_txt, err
			}

			msg_txt = fmt.Sprintf("%s reached warn limit and has been muted for 60 minutes",
		                          memberToWarn.User.FirstName)
		}
	} else {
		msg_txt = fmt.Sprintf("%s has been warned. Warns: %d/%d",
		                      memberToWarn.User.FirstName,
		                      db.Chats[chat_id].Users[user_id].Warns,
		                      db.Chats[chat_id].Config.WarnsLimit)

		err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
		if err != nil {
			msg_txt = "ERROR: can't write warn to db, check out logs"

			return msg_txt, err
		}

		return msg_txt, err
	}

	return msg_txt, err
}
