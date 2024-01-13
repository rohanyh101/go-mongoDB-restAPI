package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DevCodeChan/mongo-golang-restAPI/database"
	"github.com/DevCodeChan/mongo-golang-restAPI/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

type ServerStatusResponse struct {
	Status string `json:"status"`
}

type UserInsertedResponse struct {
	InsertedID string `json:"insertedID"`
}

type UserDeleteResponse struct {
	DeletedID string `json:"deletedID"`
}

func ServerStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := ServerStatusResponse{Status: "Server is up and running..."}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// This is not a function, this is a method, specifically a struct method...
func GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := p.ByName("id")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		log.Printf("Error retrieving user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err.Error())
		return
	}

	responseJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "%v\n", responseJSON)
	w.Write(responseJSON)
}

func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"name": user.Name})
	if err != nil {
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error occurred while checking for Name Attribute\n")
		return
	}

	if count > 0 {
		w.WriteHeader(http.StatusConflict) // Conflict status for duplicate entries
		fmt.Fprintf(w, "User with the name '%s' already exists\n", user.Name)
		return
	}

	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		log.Printf("Error inserting user: %s\n", insertErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "User not created...\n")
		return
	}

	responce := UserInsertedResponse{
		InsertedID: resultInsertionNumber.InsertedID.(primitive.ObjectID).Hex(),
	}

	responseJSON, err := json.Marshal(responce)
	if err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// fmt.Fprintf(w, "insertionId: %s\n", responseJSON)
	w.Write(responseJSON)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userId := p.ByName("id")
	user := models.User{}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		log.Printf("Error retrieving user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err.Error())
		return
	}

	// err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Err()
	// if err != nil {
	// 	log.Printf("Error retrieving user: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "Error retrieving user: %v\n", err)
	// 	return
	// }

	deleteResult, err := userCollection.DeleteOne(ctx, bson.M{"user_id": userId})
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting user: %v\n", err)
		return
	}

	if deleteResult.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found with ID: %s\n", userId)
		return
	}

	response := UserDeleteResponse{
		DeletedID: userId,
	}

	// responseJSON, err := json.Marshal(response)
	// if err != nil {
	// 	log.Printf("Error encoding JSON response: %v\n", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Deleted User ID: %s\n", responseJSON)
	// w.Write(responseJSON)
}
