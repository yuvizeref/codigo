package routes

import (
	"codigo/handlers"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PATCH")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	return r
}
