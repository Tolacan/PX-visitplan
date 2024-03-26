package main

import (
	"PX-visitplan/events"
	"PX-visitplan/models"
	"PX-visitplan/repository"
	"PX-visitplan/search"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func OnCreateVisitPlan(m events.VisitPlanMessage) {

	visit := models.VisitPlan{
		Uuid:              m.Uuid,
		Nombre:            m.Nombre,
		Ruta:              m.Ruta,
		UuidRuta:          m.UuidRuta,
		Lunes:             m.Lunes,
		Martes:            m.Martes,
		Miercoles:         m.Miercoles,
		Jueves:            m.Jueves,
		Viernes:           m.Viernes,
		Sabado:            m.Sabado,
		Domingo:           m.Domingo,
		Responsable:       m.Responsable,
		FechaRegistro:     m.FechaRegistro,
		FechaModificacion: m.FechaModificacion,
		Activo:            m.Activo,
		Clientes:          m.Clientes,
	}

	if m.TypeMsg == "create" {
		if err := search.IndexVisitPlan(context.Background(), visit); err != nil {
			log.Println(err)
		}
	}

	if m.TypeMsg == "update" {
		if err := search.UpdateVisitPlan(context.Background(), visit); err != nil {
			log.Println(err)
		}
	}
}

func listVisitPlanHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	visitsP, err := repository.ListVisitPlans(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(visitsP)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var err error
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	visitsP, err := search.SearchVisitPlan(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(visitsP)

}
