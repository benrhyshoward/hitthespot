package model

import (
	"time"

	"github.com/zmb3/spotify"
)

type User struct {
	Id      string
	Created time.Time
	Client  spotify.Client `json:"-" bson:"-"`
	Rounds  []Round        `json:"-"`
}
