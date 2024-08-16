package main

import (
	"flag"
	"log"
	tgClient "telegram-helper/clients/telegram"
	event_consumer "telegram-helper/consumer/event-consumer"
	"telegram-helper/events/telegram"
	"telegram-helper/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), files.New(storagePath))

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

func mustToken() string {
	token := flag.String("tg-token", "", "telegram bot token")
	flag.Parse()

	if *token == "" {
		log.Fatal("telegram bot token is required")
	}

	return *token
}
