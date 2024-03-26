package main

import (
	"PX-visitplan/database"
	"PX-visitplan/repository"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func newRoute() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/visitplan", createVistiPlanHandler).Methods(http.MethodPost)
	router.HandleFunc("/visitplan", updateVisitPlanHandler).Methods(http.MethodPut)
	return
}

func main() {
	err := godotenv.Load()
	var cfg Config
	err = envconfig.Process("", &cfg)

	if err != nil {
		log.Fatalf("%v", err)
	}

	psqlInfo := fmt.Sprintf("host=192.168.0.8 port=5430 user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	repo, err := database.NewPostgresRepository(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)

	//Start Server
	router := newRoute()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
