package data

import (
	"errors"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Hailemari/auth_task_manager/models"
)

var userCollection *mongo.Collection

func InitUserCollection(client *mongo.Client) {
	userCollection = client.Database("taskDB").Collection("users")
}


func CreateUser(newUser *models.User) error {
	// Check if the user already exists
	existingUser := &models.User{}
	err := userCollection.FindOne(context.TODO(), bson.M{"username": newUser.Username}).Decode(existingUser)
	if err == nil {
		return errors.New("user already exists")
	}

	// Hash the password
	err = newUser.HashPassword()
	if err != nil {
		return err
	}

	// Set the role (first user is admin, others are regular users)
	count, err := userCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	if count == 0 {
		newUser.Role = "admin"
	} else {
		newUser.Role = "user"
	}

	// Insert the new user
	_, err = userCollection.InsertOne(context.TODO(), newUser)
	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func PromoteUser(username string) error {	
	// Check the user's current role
	var user bson.M
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user not found")
		}
		return err
	}

	// Check if the user is already an admin
	if userRole, ok := user["role"].(string); ok && userRole == "admin" {
		return errors.New("user is already an admin")
	}

	// Proceed with the update if the user is not already an admin
	result, err := userCollection.UpdateOne(
		context.TODO(), 
		bson.M{"username": username},
		bson.M{"$set": bson.M{"role": "admin"}},
	)

	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}