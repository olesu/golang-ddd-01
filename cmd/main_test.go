package main

import (
	"context"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/olesu/golang-ddd-01/internal/purchase"
)

func TestRun(t *testing.T) {
	app := appContext{
		stripeService: stripeHandlerFunc(fakeStripeService),
		purchaseRepo:  fakePurchaseRepo{},
		storeRepo:     fakeStoreRepo{},
	}
	got := run(context.Background(), app)
	want := 0
	if got != want {
		t.Errorf("unexpected failure, got %d, want %d", got, want)
	}
}

type stripeHandlerFunc func(ctx context.Context, amount money.Money, cardToken string) error

func (f stripeHandlerFunc) ChargeCard(ctx context.Context, amount money.Money, cardToken string) error {
	return f(ctx, amount, cardToken)
}

func fakeStripeService(_ context.Context, _ money.Money, _ string) error {
	return nil
}

type fakePurchaseRepo struct {
}

func (r fakePurchaseRepo) Store(ctx context.Context, purchase purchase.Purchase) error { return nil }
func (r fakePurchaseRepo) Ping(ctx context.Context) error                              { return nil }

type fakeStoreRepo struct {
}

func (r fakeStoreRepo) GetStoreDiscount(ctx context.Context, storeID uuid.UUID) (float32, error) {
	return 0, nil
}
func (r fakeStoreRepo) Ping(ctx context.Context) error { return nil }
