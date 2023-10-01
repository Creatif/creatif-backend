package storage

import (
	"github.com/teris-io/shortid"
	"time"
)

var ShortId *shortid.Shortid

func initShortId() error {
	sid, err := shortid.New(1, shortid.DefaultABC, uint64(time.Now().Unix()))
	if err != nil {
		return err
	}

	ShortId = sid

	return nil
}
