package db

import (
	"context"
	"errors"
	"log"

	"github.com/benrhyshoward/hitthespot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetGuessById(userId string, roundId string, questionId string, id string) (model.Guess, error) {
	var guess model.Guess
	question, err := GetQuestionById(userId, roundId, questionId)
	if err != nil {
		return guess, err
	}

	//More likely to be requesting later guesses so searching in reverse order
	for i := len(question.Guesses) - 1; i >= 0; i-- {
		if question.Guesses[i].Id == id {
			return question.Guesses[i], nil
		}
	}
	return guess, errors.New("No guess found with id " + id + "for question" + questionId + " for round " + roundId + " for user " + userId)
}

func GetGuesses(userId string, roundId string, questionId string, filter func(model.Guess) bool) ([]model.Guess, error) {
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	question, err := GetQuestionById(user.Id, roundId, questionId)
	if err != nil {
		return nil, err
	}

	filteredGuesses := []model.Guess{}

	for _, guess := range question.Guesses {
		if filter(guess) {
			filteredGuesses = append(filteredGuesses, guess)
		}
	}
	return filteredGuesses, nil
}

func AddGuessToQuestion(userId string, roundId string, questionId string, guess model.Guess) error {
	match := bson.D{{"id", userId}}
	change := bson.D{
		{"$push", bson.D{
			{"rounds.$[r].questions.$[q].guesses", guess},
		}},
	}
	options := &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.D{
					{"r.id", roundId},
				},
				bson.D{
					{"q.id", questionId},
				},
			},
		},
	}
	updateResult, err := mongoUserCollection.UpdateOne(context.TODO(), match, change, options)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 1 {
		log.Print("Added guess " + guess.Id + " to question " + questionId)
		return nil
	}
	return errors.New("Failed to add guess " + guess.Id + " to question " + questionId)
}
