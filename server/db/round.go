package db

import (
	"context"
	"errors"
	"log"

	"github.com/benrhyshoward/hitthespot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRoundById(userId string, id string) (model.Round, error) {
	var round model.Round
	user, err := GetUserById(userId)
	if err != nil {
		return round, err
	}

	//More likely to be requesting later rounds so searching in reverse order
	for i := len(user.Rounds) - 1; i >= 0; i-- {
		if user.Rounds[i].Id == id {
			return user.Rounds[i], nil
		}
	}

	return round, errors.New("No round found with id " + id + " for user " + userId)
}

func GetRounds(userId string, filter func(model.Round) bool) ([]model.Round, error) {
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	filteredRounds := []model.Round{}
	for _, round := range user.Rounds {
		if filter(round) {
			filteredRounds = append(filteredRounds, round)
		}
	}
	return filteredRounds, nil
}

func AddRound(userId string, round model.Round) error {
	match := bson.D{{"id", userId}}
	change := bson.D{
		{"$push", bson.D{
			{"rounds", round},
		}},
	}

	updateResult, err := mongoUserCollection.UpdateOne(context.TODO(), match, change)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 1 {
		log.Print("Added round " + round.Id + " to user " + userId)
		return nil
	}
	return errors.New("Failed to add round " + round.Id + " to user " + userId)
}

func UpdateRound(userId string, round model.Round) error {
	match := bson.D{{"id", userId}}
	change := bson.D{
		{"$set", bson.D{
			{"rounds.$[r]", round},
		}},
	}
	options := &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.D{
					{"r.id", round.Id},
				},
			},
		},
	}

	_, err := mongoUserCollection.UpdateOne(context.TODO(), match, change, options)
	if err != nil {
		return err
	}

	log.Print("Updated round " + round.Id + " for user " + userId)
	return nil
}
