package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidBanUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt := ""
	var err error = nil

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: reply to user you want to ban"

		return msg_txt, err
	}

	memberToBan, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: internal error check log for the further info"

		return msg_txt, err
	}

	if memberToBan.IsAdministrator() || memberToBan.IsCreator() {
		msg_txt = "ERROR: can't ban admins"

		return msg_txt, err
	}

	var member_config tgbotapi.ChatMemberConfig
	member_config.ChatID = update.Message.ReplyToMessage.Chat.ID
	member_config.UserID = memberToBan.User.ID

	var kick_config tgbotapi.KickChatMemberConfig
	kick_config.ChatMemberConfig = member_config

	chat_id := strconv.FormatInt(member_config.ChatID, 10)
	user_id := strconv.Itoa(member_config.UserID)

	resp, err := bot.KickChatMember(kick_config)
	err = checkAPIResp(resp)
	if err != nil {
		if resp.ErrorCode == 400 {
			if resp.Description == "Bad Request: CHAT_ADMIN_REQUIRED" {
				msg_txt = "ERROR: bot is not admin"
			} else {
				msg_txt = "ERROR: " + strconv.Itoa(resp.ErrorCode) + " - " + resp.Description
			}
		} else {
			msg_txt = "ERROR: " + strconv.Itoa(resp.ErrorCode) + " - " + resp.Description
		}
		return msg_txt, err
	}

	if len(strings.Split(update.Message.Text, " ")) > 1 {
		ban_note := strings.Replace(update.Message.Text, "/ban ", "", 1)

		if db.Chats[chat_id].Users[user_id] == nil {
			db.Chats[chat_id].Users[user_id] = new(user)
		}

		db.Chats[chat_id].Users[user_id].BanNote = ban_note

		err = dbWriteChatUsers(chat_id, db.Chats[chat_id].Users)
		if err != nil {
			msg_txt = "ERROR: can't write ban description, check out logs"

			return msg_txt, err
		}
	}

	msg_txt = fmt.Sprintf("%s *banned* %s", update.Message.From.FirstName,
	                                    memberToBan.User.FirstName)

	return msg_txt, err
}
