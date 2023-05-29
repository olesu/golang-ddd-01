package main

import (
	"context"
	"log"
	"os"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	coffeeco "github.com/olesu/golang-ddd-01/internal"
	"github.com/olesu/golang-ddd-01/internal/payment"
	"github.com/olesu/golang-ddd-01/internal/purchase"
	"github.com/olesu/golang-ddd-01/internal/store"
)

type appContext struct {
	stripeService payment.StripeService
	purchaseRepo  purchase.Repository
	storeRepo     store.Repository
}

func run(ctx context.Context, app appContext) int {

	// This is a test token from Stripe's documentation. Feel free to use it, no charges will actually be made.
	cardToken := "tok_visa"

	sSvc := store.NewService(app.storeRepo)

	svc := purchase.NewService(app.stripeService, app.purchaseRepo, sSvc)

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
		log.Println(err)
		return 1
	}

	log.Println("purchase was successful")
	return 0
}

func main() {
	ctx := context.Background()

	// This is the test key from Stripe's documentation. Feel free to use it, no charges will actually be made.
	stripeTestAPIKey := "sk_test_4eC39HqLyjWDarjtT1zdp7dc"

	csvc, err := payment.NewStripeService(stripeTestAPIKey)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	// This is the credentials for mongo if you run docker-compose up in this repo.
	mongoConString := "mongodb://root:example@mongo:27017"
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

	app := appContext{
		stripeService: csvc,
		purchaseRepo:  prepo,
		storeRepo:     sRepo,
	}
	os.Exit(run(ctx, app))
}
