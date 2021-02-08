package csv

import (
	"io"
	"strings"
	"testing"

	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/itimofeev/price-store-test/internal/util"
)

func TestParser(t *testing.T) {
	testCSV := `hi;10.01
hello;7`
	reader := strings.NewReader(testCSV)
	log := logger.New()

	parser := New(log, &util.SimpleReadCloser{Reader: reader})

	next, err := parser.Next()
	require.NoError(t, err)
	require.Equal(t, "hi", next.Name)
	require.EqualValues(t, 1001, next.Price)

	next, err = parser.Next()
	require.NoError(t, err)
	require.Equal(t, "hello", next.Name)
	require.EqualValues(t, 7, next.Price)

	_, err = parser.Next()
	require.Equal(t, io.EOF, err)
}
