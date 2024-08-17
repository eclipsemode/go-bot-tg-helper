package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/eclipsemode/go-bot-tg-helper/lib/errs"
	"io"
	"time"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages found")

type Page struct {
	URL      string
	UserName string
	Created  time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", errs.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
