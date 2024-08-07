package main

import (
	"PX-visitplan/database"
	"PX-visitplan/events"
	"PX-visitplan/repository"
	"PX-visitplan/search"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresUrl          string `envconfig:"POSTGRES_URL"`
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRoute() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/visitplan", listVisitPlanHandler).Methods(http.MethodGet)
	router.HandleFunc("/visitplanheader", listVisitPlanHeaderHandler).Methods(http.MethodGet)
	router.HandleFunc("/search", searchHandler).Methods(http.MethodGet)
	return
}

func main() {
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	//Postgres DataBase
	psqlInfo := fmt.Sprintf("host=%s port=5430 user=%s password=%s dbname=%s sslmode=disable", cfg.PostgresUrl, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	repo, err := database.NewPostgresRepository(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)

	//Conexion a elasticsearch
	es, err := search.NewElastic(fmt.Sprintf("http://%s", cfg.ElasticSearchAddress))
	if err != nil {
		log.Fatal(err)
	}
	search.SetSearchRepository(es)
	defer search.Close()

	//Conexion a NATS
	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)
	defer events.Close()

	err = n.OnCreateVisitPlan(OnCreateVisitPlan)
	if err != nil {
		log.Fatalf("could not subscribe to the event: %v", err)
	}

	// Start http service
	router := newRoute()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
