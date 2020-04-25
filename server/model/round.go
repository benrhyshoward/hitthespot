package model

import "time"

type Round struct {
	Id             string
	Created        *time.Time
	Abandoned      bool
	TotalQuestions int
	Questions      []Question
}

func (r Round) InProgress() bool {
	if r.Abandoned == true {
		return false
	}
	if len(r.Questions) < r.TotalQuestions {
		return true
	}
	for _, question := range r.Questions {
		if !question.Over() {
			return true
		}
	}
	return false
}
