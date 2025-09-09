package routing

import (
	"learn_golang/internal/auth"

	"github.com/gorilla/mux"
)

func SetupAuthRoutes(router *mux.Router, authHandler *auth.AuthHandler) {
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
}
