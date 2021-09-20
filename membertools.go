package main

import (
	"errors"
	"strconv"

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

	if update.Message.ReplyToMessage.NewChatMembers != nil {
		memberConfig.UserID = (*update.Message.ReplyToMessage.NewChatMembers)[0].ID
	} else if update.Message.ReplyToMessage.LeftChatMember != nil {
		memberConfig.UserID = update.Message.ReplyToMessage.LeftChatMember.ID
	} else if update.Message.ReplyToMessage.ForwardFrom != nil {
		memberConfig.UserID = (*update.Message.ReplyToMessage.ForwardFrom).ID
	} else {
		memberConfig.UserID = update.Message.ReplyToMessage.From.ID
	}

	member, err = bot.GetChatMember(memberConfig)

	return member, err
}

func maidIsUserHasPrivileges(privilege_level int, bot *tgbotapi.BotAPI,
                             update tgbotapi.Update, db *maidDB) (bool, error) {
	var reply bool = false
	var err error  = nil

	member, err := maidGetChatMember(bot, update)
	if err != nil {
		return reply, err
	}

	chat_id := strconv.FormatInt(update.Message.Chat.ID, 10)
	user_id := strconv.Itoa(member.User.ID)

	if member.IsAdministrator() {
		reply = true

		return reply, err
	}

	if member.IsCreator() {
		reply = true

		return reply, err
	}

	if member.User.ID == config.BotAdminID {
		reply = true

		return reply, err
	}

    if db.Chats[chat_id].Users[user_id] == nil {
        db.Chats[chat_id].Users[user_id] = new(user)

    }

    // TODO: SIGSEGV null pointer dereference
	if db.Chats[chat_id].Users[user_id].Privileges >= privilege_level {
		reply = true
	}

	return reply, err
}
