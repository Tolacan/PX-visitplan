package events

import "time"

type Message interface {
	Type() string
}

type VisitPlanMessage struct {
	TypeMsg           string    `json:"typeMsg"`
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

func (msg *VisitPlanMessage) Type() string {
	return "visitplan"
}
