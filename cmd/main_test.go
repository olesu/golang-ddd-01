package main

import (
	"context"
	"errors"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/olesu/golang-ddd-01/internal/purchase"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name    string
		cardErr error
		want    int
	}{
		{
			name: "all good",
		},
		{
			name:    "failing card payment",
			cardErr: errors.New("payment failed"),
			want:    1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := appContext{
				stripeService: stripeHandlerFunc(fakeStripeService(tc.cardErr)),
				purchaseRepo:  fakePurchaseRepo{},
				storeRepo:     fakeStoreRepo{},
			}
			got := run(context.Background(), app)
			if got != tc.want {
				t.Errorf("unexpected failure, got %d, want %d", got, tc.want)
			}
		})
	}
}

type stripeHandlerFunc func(ctx context.Context, amount money.Money, cardToken string) error

func (f stripeHandlerFunc) ChargeCard(ctx context.Context, amount money.Money, cardToken string) error {
	return f(ctx, amount, cardToken)
}

func fakeStripeService(wantErr error) stripeHandlerFunc {
	return func(_ context.Context, _ money.Money, _ string) error {
		return wantErr
	}
}

type fakePurchaseRepo struct {
}

func (r fakePurchaseRepo) Store(_ context.Context, _ purchase.Purchase) error { return nil }
func (r fakePurchaseRepo) Ping(_ context.Context) error                       { return nil }

type fakeStoreRepo struct {
}

func (r fakeStoreRepo) GetStoreDiscount(_ context.Context, _ uuid.UUID) (float32, error) {
	return 0, nil
}
func (r fakeStoreRepo) Ping(_ context.Context) error { return nil }
