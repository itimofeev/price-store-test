package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/itimofeev/price-store-test/internal/model"
	"github.com/itimofeev/price-store-test/internal/pg"
	"github.com/itimofeev/price-store-test/internal/service"
	"github.com/itimofeev/price-store-test/internal/util"
)

type fakeDownloader struct {
	s string
}

func (f *fakeDownloader) GetCSV(context.Context, string) (io.ReadCloser, error) {
	return &util.SimpleReadCloser{Reader: strings.NewReader(f.s)}, nil
}

func TestHandler(t *testing.T) {
	h, productName := initApp()

	processRequest := httptest.NewRequest(http.MethodGet, "/processCSV?url=hello", nil)

	resp, err := h.Test(processRequest, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	all, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "OK", string(all))

	query := url.Values{"limit": []string{"1"}, "order": []string{"last_update DESC"}}.Encode()
	listRequest := httptest.NewRequest(http.MethodGet, "/listProducts?"+query, nil)
	resp, err = h.Test(listRequest, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	all, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	var products []model.Product
	require.NoError(t, json.Unmarshal(all, &products))

	require.Len(t, products, 1)
	require.Equal(t, productName, products[0].Name)
	require.EqualValues(t, 1001, products[0].Price)
}

func TestInvalidParams(t *testing.T) {
	h, _ := initApp()

	query := url.Values{"order": []string{"invalidColumn"}}.Encode()
	listRequest := httptest.NewRequest(http.MethodGet, "/listProducts?"+query, nil)
	resp, err := h.Test(listRequest, -1)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	all, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	const expected = "failed to list products in service: failed to select products: ERROR #42703 column \"invalidColumn\" does not exist"
	require.Equal(t, expected, string(all))
}

func initApp() (app *fiber.App, productName string) {
	productName = util.RandomID()
	testCSV := fmt.Sprintf("%s;10.01", productName)

	log := util.NewLog()
	d := &fakeDownloader{s: testCSV}
	store := pg.New(log, "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable")
	srv := service.New(log, d, store)

	app = InitApp(srv)
	return app, productName
}
