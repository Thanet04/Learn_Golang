package routing

import (
	"context"
	"encoding/json"
	"learn_golang/internal/model"
	"learn_golang/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter(userRepo *repository.UserRepository) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := userRepo.CreateUser(context.Background(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
	}).Methods("POST")

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepo.GetAllUsers(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}).Methods("GET")

	return router
}
