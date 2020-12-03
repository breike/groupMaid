package main

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidKickUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {
	msg_txt := ""
	var err error = nil

	if update.Message.ReplyToMessage == nil {
		msg_txt = "ERROR: reply to user you want to kick"

		return msg_txt, err
	}

	memberToBan, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: internal error check log for the further info"

		return msg_txt, err
	}

	if memberToBan.IsAdministrator() || memberToBan.IsCreator() {
		msg_txt = "ERROR: can't kick admins"

		return msg_txt, err
	}

	var member_config tgbotapi.ChatMemberConfig
	member_config.ChatID = update.Message.ReplyToMessage.Chat.ID
	member_config.UserID = update.Message.ReplyToMessage.From.ID

	var kick_config tgbotapi.KickChatMemberConfig
	kick_config.ChatMemberConfig = member_config
	kick_config.UntilDate = int64(update.Message.Date + 30)

	resp, err := bot.KickChatMember(kick_config)
	err = checkAPIResp(resp)
	if err != nil {
		if resp.ErrorCode == 400 {
			if resp.Description == "Bad Request: CHAT_ADMIN_REQUIRED" {
				msg_txt = "ERROR: bot is not admin"
			} else {
				msg_txt = "ERROR: " + string(resp.ErrorCode) + " - " + resp.Description
			}
		} else {
			msg_txt = "ERROR: " + string(resp.ErrorCode) + " - " + resp.Description
		}
		return msg_txt, err
	}

	msg_txt = fmt.Sprintf("%s *kicked* %s", update.Message.From.FirstName,
	                      update.Message.ReplyToMessage.From.FirstName)

	return msg_txt, err
}
