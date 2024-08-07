package models

import "time"

type VisitPlan struct {
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

type VisitPlanHeader struct {
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
}
