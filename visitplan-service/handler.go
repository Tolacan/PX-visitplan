package main

import (
	"PX-visitplan/events"
	"PX-visitplan/models"
	"PX-visitplan/repository"
	"encoding/json"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

type visitPlanRequest struct {
	Uuid              string    `json:"uuid"`
	Nombre            string    `json:"nombre"`
	Ruta              string    `json:"ruta"`
	UuidRuta          string    `json:"uuidRuta"`
	Lunes             bool      `json:"lunes"`
	Martes            bool      `json:"martes"`
	Miercoles         bool      `json:"miercoles"`
	Jueves            bool      `json:"jueves"`
	Viernes           bool      `json:"viernes"`
	Sabado            bool      `json:"sabado"`
	Domingo           bool      `json:"domingo"`
	Responsable       string    `json:"responsable"`
	FechaRegistro     time.Time `json:"fechaRegistro"`
	FechaModificacion time.Time `json:"fechaModificacion"`
	Activo            bool      `json:"activo"`
	Clientes          []string  `json:"clientes"`
}

func createVistiPlanHandler(w http.ResponseWriter, r *http.Request) {
	var req visitPlanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC()
	modifiedAt := time.Now().UTC()
	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	visitPlan := &models.VisitPlan{
		Uuid:              id.String(),
		Nombre:            req.Nombre,
		Ruta:              req.Ruta,
		UuidRuta:          req.UuidRuta,
		Lunes:             req.Lunes,
		Martes:            req.Martes,
		Miercoles:         req.Miercoles,
		Jueves:            req.Jueves,
		Viernes:           req.Viernes,
		Sabado:            req.Sabado,
		Domingo:           req.Domingo,
		Responsable:       req.Responsable,
		FechaRegistro:     createdAt,
		FechaModificacion: modifiedAt,
		Activo:            req.Activo,
		Clientes:          req.Clientes,
	}

	if err := repository.InsertVisitPlan(r.Context(), visitPlan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishCreateVisitPlan(r.Context(), visitPlan, "create"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visitPlan)
}

func updateVisitPlanHandler(w http.ResponseWriter, r *http.Request) {
	var req visitPlanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modifiedAt := time.Now().UTC()

	visitPlan := &models.VisitPlan{
		Uuid:              req.Uuid,
		Nombre:            req.Nombre,
		Ruta:              req.Ruta,
		UuidRuta:          req.UuidRuta,
		Lunes:             req.Lunes,
		Martes:            req.Martes,
		Miercoles:         req.Miercoles,
		Jueves:            req.Jueves,
		Viernes:           req.Viernes,
		Sabado:            req.Sabado,
		Domingo:           req.Domingo,
		Responsable:       req.Responsable,
		FechaRegistro:     req.FechaRegistro,
		FechaModificacion: modifiedAt,
		Activo:            req.Activo,
		Clientes:          req.Clientes,
	}

	if err := repository.UpdateVisitPlan(r.Context(), visitPlan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishCreateVisitPlan(r.Context(), visitPlan, "update"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visitPlan)

}
