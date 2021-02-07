package pg

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rs/xid"

	"github.com/itimofeev/price-store-test/cmd/internal/model"
	"github.com/itimofeev/price-store-test/cmd/internal/util"
)

func TestSaveProduct(t *testing.T) {
	ctx := context.Background()
	store := New(util.NewLog(), "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable")

	productName := RandomID()
	var product1, product2 model.Product
	var err error

	t.Run("save new product", func(t *testing.T) {
		product1, err = store.SaveProduct(ctx, time.Now(), model.ParsedProduct{
			Name:  productName,
			Price: 1001,
		})
		require.NoError(t, err)
		require.Equal(t, productName, product1.Name)
		require.EqualValues(t, 1001, product1.Price)
		require.EqualValues(t, 0, product1.UpdateCount)
	})

	t.Run("update product with same name, increase update count", func(t *testing.T) {
		product2, err = store.SaveProduct(ctx, time.Now(), model.ParsedProduct{
			Name:  productName,
			Price: 777,
		})
		require.NoError(t, err)
		require.Equal(t, productName, product2.Name)
		require.EqualValues(t, 777, product2.Price)
		require.EqualValues(t, 1, product2.UpdateCount)

		require.Condition(t, func() (success bool) {
			return product1.LastUpdate.Before(product2.LastUpdate)
		})
	})

	t.Run("list last updated product", func(t *testing.T) {
		products, err := store.ListProducts(ctx, "last_update DESC", 1, 0)
		require.NoError(t, err)
		require.Len(t, products, 1)

		require.Equal(t, product2, products[0])
	})
}

func RandomID() string {
	return xid.New().String()
}
