package main

import (
	"os"

	"github.com/naoina/toml"
)

type maidConfig struct {
	BotAdminId   int
	BotDebug     bool
	BotDirectory string
	BotWhitelist []int64
	TgBotAPI     string

}

func configInit(configPath string) (maidConfig, error) {
	var conf maidConfig
	var err error = nil

	f, err := os.Open(configPath)
	if err != nil {
		return conf, err
	}

	if err = toml.NewDecoder(f).Decode(&conf); err != nil {
		return conf, err
	}

	if err != nil {
		return conf, err
	}

	return conf, err
}
