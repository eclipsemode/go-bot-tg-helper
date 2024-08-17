package main

import (
	"context"
	"flag"
	tgClient "github.com/eclipsemode/go-bot-tg-helper/clients/telegram"
	event_consumer "github.com/eclipsemode/go-bot-tg-helper/consumer/event-consumer"
	"github.com/eclipsemode/go-bot-tg-helper/events/telegram"
	"github.com/eclipsemode/go-bot-tg-helper/storage/sqlite"
	"github.com/joho/godotenv"
	"log"
)

const (
	tgBotHost          = "api.telegram.org"
	storagePath        = "storage"
	storageSqlLitePath = storagePath + "/sqlite/data"
	batchSize          = 100
)

// 7446099344:AAEDhIsH9nV3yhoVdF0mVYCxSp6CicyAyvE

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//fStorage := files.New(storagePath)
	s, err := sqlite.New(storageSqlLitePath)
	if err != nil {
		log.Fatalf("sqlite err: %v", err)
	}

	err = s.Init(context.TODO())
	if err != nil {
		log.Fatalf("can't init sqlite db: %v", err)
	}

	eventProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), s)

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
