package store

import (
	"context"

	"github.com/google/uuid"

	coffeeco "github.com/olesu/golang-ddd-01/internal"
)

type Store struct {
	ID             uuid.UUID
	Location       string
	ProducsForSale []coffeeco.Product
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error) {
	dis, err := s.repo.GetStoreDiscount(ctx, storeID)
	if err != nil {
		return 0, err
	}
	return dis, nil
}
