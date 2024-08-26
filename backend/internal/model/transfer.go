package model

import "time"

type Transfer struct {
	ID     int64
	Type   string
	From   string
	To     string
	Amount float64
	Time   time.Time
}
