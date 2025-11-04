package config

import (
	"errors"
	"os"
	"tg-app/model"
)

func NewConfig() (*model.Cfg, error) {
	config := model.Cfg{
		BotAPI: os.Getenv("TOKEN_BOT"),
		DBConnect: os.Getenv("DB_CONNECT"),
		BotPass: os.Getenv("PASS_BOT"),
	}

	if config.BotAPI == "" || config.DBConnect == "" || config.BotPass == "" {
		return nil, errors.New("error config")
	}

	return &config, nil
}