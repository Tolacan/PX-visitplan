package events

import (
	"PX-visitplan/models"
	"bytes"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn *nats.Conn
	//Connecion para cuando un evento ha sido creado
	visitplanCreatedSub *nats.Subscription
	//Canal para cuando un evento ha sido creado
	visitplanCreatedChan chan VisitPlanMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	return &NatsEventStore{
		conn: conn,
	}, nil
}

func (store *NatsEventStore) Close() {
	if store.conn != nil {
		store.conn.Close()
	}
	if store.visitplanCreatedSub != nil {
		store.visitplanCreatedSub.Unsubscribe()
	}
	close(store.visitplanCreatedChan)
}

func (store *NatsEventStore) encodeMessage(m Message) ([]byte, error) {

	//Se crea un buffer nuevo de tipo bytes
	b := bytes.Buffer{}
	b.Reset()
	//Codifica el mensaje en el buffer de Message a bytes
	err := json.NewEncoder(&b).Encode(m)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (store *NatsEventStore) PublishCreateVisitPlan(ctx context.Context, visit *models.VisitPlan, tipo string) error {

	msg := &VisitPlanMessage{
		TypeMsg:           tipo,
		Uuid:              visit.Uuid,
		Nombre:            visit.Nombre,
		Ruta:              visit.Ruta,
		UuidRuta:          visit.UuidRuta,
		Lunes:             visit.Lunes,
		Martes:            visit.Martes,
		Miercoles:         visit.Miercoles,
		Jueves:            visit.Jueves,
		Viernes:           visit.Viernes,
		Sabado:            visit.Sabado,
		Domingo:           visit.Domingo,
		Responsable:       visit.Responsable,
		FechaRegistro:     visit.FechaRegistro,
		FechaModificacion: visit.FechaModificacion,
		Activo:            visit.Activo,
		Clientes:          visit.Clientes,
	}

	//Codifica el mensaje
	b, err := store.encodeMessage(msg)
	if err != nil {
		return err
	}

	//Publica el mensaje
	return store.conn.Publish(msg.Type(), b)

}

func (store *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Reset()
	b.Write(data)
	return json.NewDecoder(&b).Decode(m)
}

func (store *NatsEventStore) OnCreateVisitPlan(f func(VisitPlanMessage)) (err error) {
	msg := VisitPlanMessage{}
	store.visitplanCreatedSub, err = store.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		store.decodeMessage(m.Data, &msg)
		f(msg)
	})
	return
}

func (store *NatsEventStore) SubscribeCreateVisitPlan(ctx context.Context) (<-chan VisitPlanMessage, error) {
	msg := VisitPlanMessage{}
	store.visitplanCreatedChan = make(chan VisitPlanMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	store.visitplanCreatedSub, err = store.conn.ChanSubscribe(msg.Type(), ch)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case m := <-ch:
				store.decodeMessage(m.Data, &msg)
				store.visitplanCreatedChan <- msg
			}
		}
	}()
	return (<-chan VisitPlanMessage)(store.visitplanCreatedChan), nil
}
