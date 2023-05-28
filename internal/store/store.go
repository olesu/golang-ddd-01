package store

import (
	"github.com/google/uuid"

	coffeeco "github.com/olesu/golang-ddd-01/internal"
)

type Store struct {
	ID             uuid.UUID
	Location       string
	ProducsForSale []coffeeco.Product
}
