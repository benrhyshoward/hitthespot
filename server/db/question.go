package db

import (
	"context"
	"errors"
	"log"

	"github.com/benrhyshoward/hitthespot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetQuestionById(userId string, roundId string, id string) (model.Question, error) {
	var question model.Question
	round, err := GetRoundById(userId, roundId)
	if err != nil {
		return question, err
	}

	//More likely to be requesting later questions so searching in reverse order
	for i := len(round.Questions) - 1; i >= 0; i-- {
		if round.Questions[i].Id == id {
			return round.Questions[i], nil
		}
	}

	return question, errors.New("No question found with id " + id + " for round " + roundId + " for user " + userId)
}

func GetQuestions(userId string, roundId string, filter func(model.Question) bool) ([]model.Question, error) {
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	round, err := GetRoundById(user.Id, roundId)
	if err != nil {
		return nil, err
	}

	filteredQuestions := []model.Question{}

	for _, question := range round.Questions {
		if filter(question) {
			filteredQuestions = append(filteredQuestions, question)
		}
	}
	return filteredQuestions, nil
}

func AddQuestionToRound(userId string, roundId string, question model.Question) error {
	match := bson.D{{"id", userId}}
	change := bson.D{
		{"$push", bson.D{
			{"rounds.$[r].questions", question},
		}},
	}
	options := &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.D{
					{"r.id", roundId},
				},
			},
		},
	}
	updateResult, err := mongoUserCollection.UpdateOne(context.TODO(), match, change, options)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 1 {
		log.Print("Added question " + question.Id + " to round " + roundId)
		return nil
	}
	return errors.New("Failed to add question " + question.Id + " to round " + roundId)
}

func UpdateQuestion(userId string, roundId string, question model.Question) error {
	match := bson.D{{"id", userId}}
	change := bson.D{
		{"$set", bson.D{
			{"rounds.$[r].questions.$[q]", question},
		}},
	}
	options := &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.D{
					{"r.id", roundId},
				},
				bson.D{
					{"q.id", question.Id},
				},
			},
		},
	}

	_, err := mongoUserCollection.UpdateOne(context.TODO(), match, change, options)
	if err != nil {
		return err
	}

	log.Print("Updated question " + question.Id + " for round " + roundId)
	return nil
}
