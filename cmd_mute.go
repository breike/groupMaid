package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidMuteUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {
	msg                   := ""
	var err error          = nil

	var mute_time            int

	canSendMessages       := new(bool)
	canSendMediaMessages  := new(bool)
	canSendOtherMessages  := new(bool)
	canAddWebPagePreviews := new(bool)
	post_msg              := update.Message

	memberFromCmd, err := maidGetChatMember(bot, update)
	if err != nil {
		msg = "ERROR: internal error check log for the further info"

		return msg, err
	}

	if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
		msg = "ERROR: not admin"

		return msg, err
	}

	if post_msg.ReplyToMessage == nil {
		msg = "ERROR: reply to user you want to ban"

		return msg, err
	}

	memberToBan, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg = "ERROR: internal error check log for the further info"

		return msg, err
	}

	if memberToBan.IsAdministrator() || memberToBan.IsCreator() {
		msg = "ERROR: can't mute admins"

		return msg, err
	}

	if len(strings.Fields(post_msg.Text)) > 1 {
		mute_time, err = strconv.Atoi(strings.Fields(post_msg.Text)[1])
		if err != nil {
			msg = "ERROR: looks like that entered time is not an integer"
			return msg, err
		}
	} else {
		mute_time = 60
	}

	var member_config tgbotapi.ChatMemberConfig
	member_config.ChatID = post_msg.ReplyToMessage.Chat.ID
	member_config.UserID = post_msg.ReplyToMessage.From.ID

	*canSendMessages       = false
	*canSendMediaMessages  = false
	*canSendOtherMessages  = false
	*canAddWebPagePreviews = false

	var mute_config tgbotapi.RestrictChatMemberConfig
	mute_config.UntilDate             = int64(post_msg.Date + (mute_time * 60))
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
				msg = "ERROR: bot is not admin"
			} else {
				msg = "ERROR: " + string(resp.ErrorCode) + " - " + resp.Description
			}
		} else {
			msg = "ERROR: " + string(resp.ErrorCode) + " - " + resp.Description
		}
		return msg, err
	}

	msg = fmt.Sprintf("%s *muted* %s for %d minutes", post_msg.From.FirstName,
	                  post_msg.ReplyToMessage.From.FirstName, mute_time)
	return msg, err
}

func maidUnmuteUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) (string, error) {
	msg                   := ""
	var err error          = nil

	canSendMessages       := new(bool)
	canSendMediaMessages  := new(bool)
	canSendOtherMessages  := new(bool)
	canAddWebPagePreviews := new(bool)
	post_msg              := update.Message

	memberFromCmd, err := maidGetChatMember(bot, update)
	if err != nil {
		msg = "ERROR: internal error check log for the further info"

		return msg, err
	}

	if !(memberFromCmd.IsAdministrator()) && !(memberFromCmd.IsCreator()) {
		msg = "ERROR: not admin"

		return msg, err
	}

	if post_msg.ReplyToMessage == nil {
		msg = "ERROR: reply to user you want to ban"

		return msg, err
	}

	memberToBan, err := maidGetReplyChatMember(bot, update)
	if err != nil {
		msg = "ERROR: internal error check log for the further info"

		return msg, err
	}

	if memberToBan.IsAdministrator() || memberToBan.IsCreator() {
		msg = "ERROR: can't mute admins"

		return msg, err
	}

	var member_config tgbotapi.ChatMemberConfig
	member_config.ChatID = post_msg.ReplyToMessage.Chat.ID
	member_config.UserID = post_msg.ReplyToMessage.From.ID

	*canSendMessages       = true
	*canSendMediaMessages  = true
	*canSendOtherMessages  = true
	*canAddWebPagePreviews = true

	var mute_config tgbotapi.RestrictChatMemberConfig
	mute_config.UntilDate             = int64(post_msg.Date)
	mute_config.ChatMemberConfig      = member_config
	mute_config.CanSendMessages       = canSendMessages
	mute_config.CanSendMediaMessages  = canSendMediaMessages
	mute_config.CanSendOtherMessages  = canSendOtherMessages
	mute_config.CanAddWebPagePreviews = canAddWebPagePreviews

	resp, err := bot.RestrictChatMember(mute_config)
	err = checkAPIResp(resp)
	if msg != "" {
		return msg, err
	}

	msg = fmt.Sprintf("%s *unmuted*",
	                  post_msg.ReplyToMessage.From.FirstName)
	return msg, err
}
