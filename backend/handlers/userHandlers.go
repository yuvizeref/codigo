package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"codigo/models"
	"codigo/utils"

	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := utils.GetUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := utils.GetUser(params["id"])
	if err != nil {
		log.Printf("Error fetching user with ID %s: %v", params["id"], err)
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if (user == models.User{}) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := utils.CreateUser(newUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedUser models.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err = utils.UpdateUser(params["id"], updatedUser)
	if err != nil {
		log.Printf("Error updating user with ID %s: %v", params["id"], err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	err := utils.DeleteUser(params["id"])
	if err != nil {
		log.Printf("Error deleting user with ID %s: %v", params["id"], err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
