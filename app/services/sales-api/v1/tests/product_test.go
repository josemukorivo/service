package tests

import (
	"context"
	"os"
	"runtime/debug"
	"testing"

	"github.com/ardanlabs/service/app/services/sales-api/v1/build/all"
	"github.com/ardanlabs/service/business/data/dbtest"
	"github.com/ardanlabs/service/business/web/v1/mux"
)

func Test_Product(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, c, "Test_Product")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	app := appTest{
		Handler: mux.WebAPI(mux.Config{
			Shutdown: make(chan os.Signal, 1),
			Log:      dbTest.Log,
			Auth:     dbTest.V1.Auth,
			DB:       dbTest.DB,
		}, all.Routes()),
		userToken:  dbTest.TokenV1("user@example.com", "gophers"),
		adminToken: dbTest.TokenV1("admin@example.com", "gophers"),
	}

	// -------------------------------------------------------------------------

	sd, err := createProductSeed(context.Background(), dbTest)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	app.test(t, productQuery200(t, app, sd), "product-query-200")
	app.test(t, productQueryByID200(t, app, sd), "product-querybyid-200")

	app.test(t, productCreate200(t, app, sd), "product-create-200")
	app.test(t, productCreate401(t, app, sd), "product-create-401")
	app.test(t, productCreate400(t, app, sd), "product-create-400")

	app.test(t, productUpdate200(t, app, sd), "product-update-200")
	app.test(t, productUpdate401(t, app, sd), "product-update-401")
	app.test(t, productUpdate400(t, app, sd), "product-update-400")

	app.test(t, productDelete200(t, app, sd), "product-delete-200")
	app.test(t, productDelete401(t, app, sd), "product-delete-401")
}
