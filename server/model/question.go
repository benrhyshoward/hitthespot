package model

import (
	"encoding/json"
	"time"
)

const (
	FreeText       = "FreeText"
	MultipleChoice = "MultipleChoice"
	Portmanteau    = "Portmanteau"
)

type Question struct {
	Id          string
	Type        string
	Created     *time.Time
	Abandoned   bool
	AbandonedAt *time.Time
	Description string
	Content     string
	Images      []string
	Options     []string
	Answer      Answer
	Guesses     []Guess
}

type Answer struct {
	Value     string
	ExtraInfo string
}

func (q Question) Over() bool {
	if q.Abandoned {
		return true
	}
	for _, guess := range q.Guesses {
		if guess.Correct {
			return true
		}
	}
	return false
}

//Only want to marshal the answer if question is answered or abandoned
func (q Question) MarshalJSON() ([]byte, error) {
	if !q.Over() {
		q.Answer = Answer{}
	}
	//Can't just call json.Marshal(q) as this creates an infinte loop, so copying each field
	j, err := json.Marshal(struct {
		Id          string
		Type        string
		Created     *time.Time
		Abandoned   bool
		AbandonedAt *time.Time
		Description string
		Content     string
		Images      []string
		Options     []string
		Answer      Answer
		Guesses     []Guess
	}{
		Id:          q.Id,
		Type:        q.Type,
		Created:     q.Created,
		Abandoned:   q.Abandoned,
		AbandonedAt: q.AbandonedAt,
		Description: q.Description,
		Content:     q.Content,
		Images:      q.Images,
		Options:     q.Options,
		Answer:      q.Answer,
		Guesses:     q.Guesses,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}
