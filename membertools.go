package main

import (
	"errors"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidGetChatMember(bot *tgbotapi.BotAPI, update tgbotapi.Update) (tgbotapi.ChatMember, error) {
	var member tgbotapi.ChatMember
	var err error = nil

	if update.Message == nil {
		err = errors.New("can't get message")
		return member, err
	}

	var memberConfig tgbotapi.ChatConfigWithUser
	memberConfig.ChatID = update.Message.Chat.ID
	memberConfig.UserID = update.Message.From.ID

	member, err = bot.GetChatMember(memberConfig)

	return member, err
}

func maidGetReplyChatMember(bot *tgbotapi.BotAPI, update tgbotapi.Update) (tgbotapi.ChatMember, error) {
	var member tgbotapi.ChatMember
	var err error = nil

	if update.Message.ReplyToMessage == nil {
		err = errors.New("can't get reply from message")
		return member, err
	}

	var memberConfig tgbotapi.ChatConfigWithUser
	memberConfig.ChatID = update.Message.ReplyToMessage.Chat.ID
	memberConfig.UserID = update.Message.ReplyToMessage.From.ID

	member, err = bot.GetChatMember(memberConfig)

	return member, err
}
