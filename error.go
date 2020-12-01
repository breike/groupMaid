package main

import (
	"errors"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func checkAPIResp(resp tgbotapi.APIResponse) (error) {
	var err error = nil

	errStr := fmt.Sprintf("API Error: %d - %s", resp.ErrorCode,
	                                            resp.Description)

	if resp.Ok {
		return err
	}

	err = errors.New(errStr)

	return err
}
