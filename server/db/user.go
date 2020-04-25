package db

import (
	"context"
	"log"

	"github.com/benrhyshoward/hitthespot/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserById(id string) (model.User, error) {
	var user model.User
	filter := bson.D{{"id", id}}
	err := mongoUserCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUsers() ([]model.User, error) {
	var users []model.User
	cur, err := mongoUserCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {

		var elem model.User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		users = append(users, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())

	return users, nil
}

func AddUser(user model.User) error {
	_, err := mongoUserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	log.Print("Added user " + user.Id)
	return nil
}
