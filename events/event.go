package events

import (
	"PX-visitplan/models"
	"context"
)

type EventStore interface {
	Close()
	PublishCreateVisitPlan(ctx context.Context, visit *models.VisitPlan, tipo string) error
	SubscribeCreateVisitPlan(ctx context.Context) (<-chan VisitPlanMessage, error)
	OnCreateVisitPlan(f func(VisitPlanMessage)) error
}

var eventStore EventStore

func SetEventStore(store EventStore) {
	eventStore = store
}

func Close() {
	eventStore.Close()
}

func PublishCreateVisitPlan(ctx context.Context, visit *models.VisitPlan, tipo string) error {
	return eventStore.PublishCreateVisitPlan(ctx, visit, tipo)
}

func SubscribeCreateVisitPlan(ctx context.Context) (<-chan VisitPlanMessage, error) {
	return eventStore.SubscribeCreateVisitPlan(ctx)
}

func OnCreateVisitPlan(f func(message VisitPlanMessage)) error {
	return eventStore.OnCreateVisitPlan(f)
}
