package main

import (
	"context"
	"log"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	coffeeco "github.com/olesu/golang-ddd-01/internal"
	"github.com/olesu/golang-ddd-01/internal/payment"
	"github.com/olesu/golang-ddd-01/internal/purchase"
	"github.com/olesu/golang-ddd-01/internal/store"
)

func main() {

	ctx := context.Background()

	// This is the test key from Stripe's documentation. Feel free to use it, no charges will actually be made.
	stripeTestAPIKey := "sk_test_4eC39HqLyjWDarjtT1zdp7dc"

	// This is a test token from Stripe's documentation. Feel free to use it, no charges will actually be made.
	cardToken := "tok_visa"

	// This is the credentials for mongo if you run docker-compose up in this repo.
	mongoConString := "mongodb://root:example@mongo:27017"
	csvc, err := payment.NewStripeService(stripeTestAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	prepo, err := purchase.NewMongoRepo(ctx, mongoConString)
	if err != nil {
		log.Fatal(err)
	}
	if err := prepo.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	sRepo, err := store.NewMongoRepo(ctx, mongoConString)
	if err != nil {
		log.Fatal(err)
	}
	if err := sRepo.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	sSvc := store.NewService(sRepo)

	svc := purchase.NewService(csvc, prepo, sSvc)

	someStoreID := uuid.New()

	pur := &purchase.Purchase{
		CardToken: &cardToken,
		Store: store.Store{
			ID: someStoreID,
		},
		ProductsToPurchase: []coffeeco.Product{{
			ItemName:  "item1",
			BasePrice: *money.New(3300, "USD"),
		}},
		PaymentMeans: payment.MEANS_CARD,
	}
	if err := svc.CompletePurchase(ctx, pur, nil, someStoreID); err != nil {
		log.Fatal(err)
	}

	log.Println("purchase was successful")
}
