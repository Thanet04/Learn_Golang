package routing

import (
	"context"
	"encoding/json"
	"learn_golang/internal/auth"
	"learn_golang/internal/model"
	"learn_golang/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func SetupRouter(userRepo *repository.UserRepository) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/users/{id}", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepo.GetUser(context.Background(), mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}))).Methods("GET")

	router.Handle("/users/{id}", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]

		result, err := userRepo.DeleteUser(context.Background(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User deleted successfully"))
	}))).Methods("DELETE")

	router.Handle("/users/{id}", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := mux.Vars(r)["id"]

		result, err := userRepo.UpdateUser(context.Background(), userID, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User updated successfully"))
	}))).Methods("PUT")

	router.Handle("/users", auth.JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepo.GetAllUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var response []UserResponse
		for _, user := range users {
			response = append(response, UserResponse{
				Name:  user.Name,
				Email: user.Email,
				Age:   user.Age,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))).Methods("GET")

	return router
}
