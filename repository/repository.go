package repository

import (
	"PX-visitplan/models"
	"context"
)

type Repository interface {
	Close()
	InsertVisitPlan(ctx context.Context, visit *models.VisitPlan) error
	UpdateVisitPlan(ctx context.Context, visit *models.VisitPlan) error
	ListVisitPlans(ctx context.Context) ([]*models.VisitPlan, error)
	ListVisitPlansHeader(ctx context.Context) ([]*models.VisitPlanHeader, error)
}

var repo Repository

func SetRepository(r Repository) {
	repo = r
}

func Close() {
	repo.Close()
}

func InsertVisitPlan(ctx context.Context, visit *models.VisitPlan) error {
	return repo.InsertVisitPlan(ctx, visit)
}

func ListVisitPlans(ctx context.Context) ([]*models.VisitPlan, error) {
	return repo.ListVisitPlans(ctx)
}

func ListVisitPlansHeader(ctx context.Context) ([]*models.VisitPlanHeader, error) {
	return repo.ListVisitPlansHeader(ctx)
}

func UpdateVisitPlan(ctx context.Context, visit *models.VisitPlan) error {
	return repo.UpdateVisitPlan(ctx, visit)
}
