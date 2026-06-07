package repository

import (
	"context"

	"lineoa-miniapp/domain"
)

type EventRepository interface {
	FindByID(ctx context.Context, eventID uint) (*domain.EventInfo, error)
	FindAll(ctx context.Context) ([]domain.EventInfo, error)
}
