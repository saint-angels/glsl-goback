package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Artwork struct {
	ID int
	Created time.Time
}
