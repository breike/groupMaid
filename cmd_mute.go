package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidMuteUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {
	msg_txt               := ""
	var err error          = nil

	var mute_time            int

	canSendMessages       := new(bool)
	canSendMediaMessages  := new(bool)
	canSendOtherMessages  := new(bool)
	canAddWebPagePreviews := new(bool)

	memberFromCmd, err := maidGetChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: internal error check log for the further info"

		return msg_txt, err
	}

	if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
		msg_txt = "ERROR: not admin"

		return msg_txt, err
	}

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
		msg_txt = "ERROR: can't mute admins"

		return msg_txt, err
	}

	if len(strings.Fields(update.Message.Text)) > 1 {
		mute_time, err = strconv.Atoi(strings.Fields(update.Message.Text)[1])
		if err != nil {
			msg_txt = "ERROR: looks like that entered time is not an integer"
			return msg_txt, err
		}
	} else {
		mute_time = 60
	}

	var member_config tgbotapi.ChatMemberConfig
	member_config.ChatID = update.Message.ReplyToMessage.Chat.ID
	member_config.UserID = update.Message.ReplyToMessage.From.ID

	*canSendMessages       = false
	*canSendMediaMessages  = false
	*canSendOtherMessages  = false
	*canAddWebPagePreviews = false

	var mute_config tgbotapi.RestrictChatMemberConfig
	mute_config.UntilDate             = int64(update.Message.Date + (mute_time * 60))
	mute_config.ChatMemberConfig      = member_config
	mute_config.CanSendMessages       = canSendMessages
	mute_config.CanSendMediaMessages  = canSendMediaMessages
	mute_config.CanSendOtherMessages  = canSendOtherMessages
	mute_config.CanAddWebPagePreviews = canAddWebPagePreviews

	resp, err := bot.RestrictChatMember(mute_config)
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

	msg_txt = fmt.Sprintf("%s *muted* %s for %d minutes", update.Message.From.FirstName,
	                  update.Message.ReplyToMessage.From.FirstName, mute_time)
	return msg_txt, err
}

func maidUnmuteUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {
	msg_txt               := ""
	var err error          = nil

	canSendMessages       := new(bool)
	canSendMediaMessages  := new(bool)
	canSendOtherMessages  := new(bool)
	canAddWebPagePreviews := new(bool)

	memberFromCmd, err := maidGetChatMember(bot, update)
	if err != nil {
		msg_txt = "ERROR: internal error check log for the further info"

		return msg_txt, err
	}

	if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
		msg_txt = "ERROR: not admin"

		return msg_txt, err
	}

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
		msg_txt = "ERROR: can't mute admins"

		return msg_txt, err
	}

	var member_config tgbotapi.ChatMemberConfig
	member_config.ChatID = update.Message.ReplyToMessage.Chat.ID
	member_config.UserID = update.Message.ReplyToMessage.From.ID

	*canSendMessages       = true
	*canSendMediaMessages  = true
	*canSendOtherMessages  = true
	*canAddWebPagePreviews = true

	var mute_config tgbotapi.RestrictChatMemberConfig
	mute_config.UntilDate             = int64(update.Message.Date)
	mute_config.ChatMemberConfig      = member_config
	mute_config.CanSendMessages       = canSendMessages
	mute_config.CanSendMediaMessages  = canSendMediaMessages
	mute_config.CanSendOtherMessages  = canSendOtherMessages
	mute_config.CanAddWebPagePreviews = canAddWebPagePreviews

	resp, err := bot.RestrictChatMember(mute_config)
	err = checkAPIResp(resp)
	if msg_txt != "" {
		return msg_txt, err
	}

	msg_txt = fmt.Sprintf("%s *unmuted*",
	                  update.Message.ReplyToMessage.From.FirstName)
	return msg_txt, err
}
