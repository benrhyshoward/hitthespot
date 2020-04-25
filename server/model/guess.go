package model

import "time"

type Guess struct {
	Id      string
	Created *time.Time
	Content string
	Correct bool
	Score   int
}
