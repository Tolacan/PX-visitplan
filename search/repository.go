package search

import (
	"PX-visitplan/models"
	"context"
)

type SearchRepository interface {
	Close()
	IndexVisitPlan(ctx context.Context, visitPlan models.VisitPlan) error
	UpdateVisitPlan(ctx context.Context, visitPlan models.VisitPlan) error
	SearchVisitPlan(ctx context.Context, query string) ([]models.VisitPlan, error)
}

var repo SearchRepository

func SetSearchRepository(r SearchRepository) {
	repo = r
}

func Close() {
	repo.Close()
}

func IndexVisitPlan(ctx context.Context, visitPlan models.VisitPlan) error {
	return repo.IndexVisitPlan(ctx, visitPlan)
}

func UpdateVisitPlan(ctx context.Context, visitPlan models.VisitPlan) error {
	return repo.UpdateVisitPlan(ctx, visitPlan)
}

func SearchVisitPlan(ctx context.Context, query string) ([]models.VisitPlan, error) {
	return repo.SearchVisitPlan(ctx, query)
}
