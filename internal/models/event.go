package models

import "time"

type Event struct {
	ID        int64
	Message   string
	Timestamp time.Time
}
