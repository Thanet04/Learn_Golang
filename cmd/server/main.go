package main

import (
	"learn_golang/config"
	"learn_golang/internal/repository"
	"learn_golang/internal/routing"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	client := config.ConnectToMongo(cfg)
	defer config.DisconnectMongo(client)

	userRepo := repository.NewUserRepository(client, cfg.DBName)
	router := routing.SetupRouter(userRepo)

	log.Fatal(Run(router, cfg.ServerPort))
}

func Run(handler http.Handler, port string) error {
	addr := ":" + port
	log.Printf("Starting server on %s...\n", addr)
	return http.ListenAndServe(addr, handler)
}
