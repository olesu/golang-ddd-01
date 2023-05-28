package purchase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	coffeeco "github.com/olesu/golang-ddd-01/internal"
	"github.com/olesu/golang-ddd-01/internal/loyalty"
	"github.com/olesu/golang-ddd-01/internal/payment"
	"github.com/olesu/golang-ddd-01/internal/store"
)

type CardChargeService interface {
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}

type StoreService interface {
	GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error)
}

type Purchase struct {
	id                 uuid.UUID
	Store              store.Store
	ProductsToPurchase []coffeeco.Product
	total              money.Money
	PaymentMeans       payment.Means
	timeOfPurchase     time.Time
	cardToken          *string
}

func (p *Purchase) validateAndEnrich() error {
	if len(p.ProductsToPurchase) == 0 {
		return errors.New("purchase must consist of at least one product")
	}
	p.total = *money.New(0, "USD")
	for _, v := range p.ProductsToPurchase {
		newTotal, _ := p.total.Add(&v.BasePrice)
		p.total = *newTotal
	}
	if p.total.IsZero() {
		return errors.New("likely mistake; purchase should never be 0. Please validate")
	}
	p.id = uuid.New()
	p.timeOfPurchase = time.Now()
	return nil
}

type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
	storeService StoreService
	repo         store.Repository
}

func (s Service) CompletePurchase(ctx context.Context, purchase *Purchase, coffeeBuxCard *loyalty.CoffeeBux, storeID uuid.UUID) error {
	if err := purchase.validateAndEnrich(); err != nil {
		return err
	}
	if err := s.calculateStoreSpecificDiscount(ctx, storeID, purchase); err != nil {
		return err
	}
	switch purchase.PaymentMeans {
	case payment.MEANS_CARD:
		if err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.cardToken); err != nil {
			return errors.New("card charge failed, cancelling purchase")
		}
	case payment.MEANS_CASH:
		//
	case payment.MEANS_COFFEEBUX:
		if err := coffeeBuxCard.Pay(ctx, purchase.ProductsToPurchase); err != nil {
			return fmt.Errorf("failed to charge loyalty card: %w", err)
		}
	default:
		return errors.New("unknown payment type")
	}
	if err := s.purchaseRepo.Store(ctx, *purchase); err != nil {
		return errors.New("failed to store purchase")
	}
	if coffeeBuxCard != nil {
		coffeeBuxCard.AddStamp()
	}
	return nil
}

func (s Service) calculateStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID, purchase *Purchase) error {
	discount, err := s.storeService.GetStoreSpecificDiscount(ctx, storeID)
	if err != nil && err != store.ErrNoDiscount {
		return fmt.Errorf("failed to get discount: %w", err)
	}
	purchasePrice := purchase.total
	if discount > 0 {
		purchase.total = *purchasePrice.Multiply(int64(100 - discount))
	}
	return nil
}

func (s Service) GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error) {
	dis, err := s.repo.GetStoreDiscount(ctx, storeID)
	if err != nil {
		return 0, err
	}
	return dis, nil
}