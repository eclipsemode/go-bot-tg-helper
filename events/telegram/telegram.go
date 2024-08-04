package telegram

import "telegram-helper/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
}

func New(client *telegram.Client) {

}
