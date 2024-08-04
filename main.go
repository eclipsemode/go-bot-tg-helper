package main

import (
	"flag"
	"log"
	"telegram-helper/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(mustToken(), tgBotHost)
}

func mustToken() string {
	token := flag.String("tg-token", "", "telegram bot token")
	flag.Parse()

	if *token == "" {
		log.Fatal("telegram bot token is required")
	}

	return *token
}
