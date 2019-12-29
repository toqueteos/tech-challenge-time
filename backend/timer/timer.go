package timer

import (
	"errors"
	"pento-tech-challenge/backend/null"
)

var (
	ErrTimerInvalidID = errors.New("timer: invalid ID")
)

type Timer struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	DateStart null.Time `db:"date_start" json:"date_start,omitempty"`
	DateEnd   null.Time `db:"date_end" json:"date_end,omitempty"`
}
