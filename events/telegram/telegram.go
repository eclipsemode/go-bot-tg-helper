package telegram

import (
	"errors"
	"telegram-helper/clients/telegram"
	"telegram-helper/events"
	"telegram-helper/lib/errs"
	"telegram-helper/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errs.Wrap("can't get events", err)
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, u := range update {
		res = append(res, event(u))
	}

	p.offset = update[len(update)-1].ID + 1

	return res, nil
}

func (p Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		p.proccessMessage(event)
	default:
		return errs.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) proccessMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errs.Wrap("can't process message", err)
	}
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, errs.Wrap("can't get meta from event", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
