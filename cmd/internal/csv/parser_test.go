package csv

import (
	"io"
	"strings"
	"testing"

	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/itimofeev/price-store-test/cmd/internal/util"
)

func TestName(t *testing.T) {
	testCSV := `hi;there
hello;world`
	reader := strings.NewReader(testCSV)
	log := logger.New()

	parser := New(log, &util.SimpleReadCloser{Reader: reader})

	next, err := parser.Next()
	require.NoError(t, err)
	require.Equal(t, "hi", next.Name)
	require.Equal(t, "there", next.Price)

	next, err = parser.Next()
	require.NoError(t, err)
	require.Equal(t, "hello", next.Name)
	require.Equal(t, "world", next.Price)

	_, err = parser.Next()
	require.Equal(t, io.EOF, err)
}
