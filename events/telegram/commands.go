package telegram

import (
	"context"
	"errors"
	"github.com/eclipsemode/go-bot-tg-helper/clients/telegram"
	"github.com/eclipsemode/go-bot-tg-helper/lib/errs"
	"github.com/eclipsemode/go-bot-tg-helper/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(context.Background(), chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(context.Background(), chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = errs.WrapIfErr("can't do command: save page", err)
	}()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() {
		err = errs.WrapIfErr("can't do command: send random", err)
	}()

	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(text string) error {
		return tg.SendMessage(chatID, text)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
