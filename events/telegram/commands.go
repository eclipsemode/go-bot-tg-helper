package telegram

import (
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/cmd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {

	}

	switch text {
	case RndCmd:
	case HelpCmd:
	case StartCmd:
	default:

	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
