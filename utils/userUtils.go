package utils

import (
	"codigo/db"
	"codigo/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers() ([]models.User, error) {
	collection := db.DB.Collection("users")

	filter := bson.D{{}}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Printf("Error decoding user: %v", err)
			return nil, err
		}
		user.Password = ""
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return users, nil
}

func GetUser(userID string) (models.User, error) {
	collection := db.DB.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid userID format: %v", err)
		return models.User{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	var user models.User

	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, nil
	}
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return user, err
	}

	user.Password = ""
	return user, nil
}

func CreateUser(user models.User) (models.User, error) {
	collection := db.DB.Collection("users")

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return models.User{}, err
	}

	user.Password = string(hashedPassword)

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return models.User{}, err
	}

	user.Password = ""

	return user, nil
}

func UpdateUser(userID string, user models.User) (models.User, error) {
	collection := db.DB.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid userID format: %v", err)
		return models.User{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	updateFields := bson.D{}

	if user.Username != "" {
		updateFields = append(updateFields, bson.E{Key: "username", Value: user.Username})
	}
	if user.Email != "" {
		updateFields = append(updateFields, bson.E{Key: "email", Value: user.Email})
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return models.User{}, err
		}
		updateFields = append(updateFields, bson.E{Key: "password", Value: string(hashedPassword)})
	}
	if user.Name != "" {
		updateFields = append(updateFields, bson.E{Key: "name", Value: user.Name})
	}
	if user.Admin {
		updateFields = append(updateFields, bson.E{Key: "admin", Value: user.Admin})
	}

	if len(updateFields) == 0 {
		return models.User{}, nil
	}

	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return models.User{}, err
	}

	user.Password = ""

	return user, nil
}

func DeleteUser(userID string) error {
	collection := db.DB.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid userID format: %v", err)
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	return nil
}
