package event_consumer

import (
	"log"
	"telegram-helper/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)

		if err != nil {
			log.Printf("[ERROR] Failed to fetch events: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvent(gotEvents); err != nil {
			log.Print(err)

			continue
		}

	}
}

func (c *Consumer) handleEvent(event []events.Event) error {
	for _, event := range event {
		log.Printf("got event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("cannot process event: %s", err.Error())

			continue
		}
	}

	return nil
}
