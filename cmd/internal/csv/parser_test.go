package csv

import (
	"io"
	"strings"
	"testing"

	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	testCSV := `hi;there
hello;world`
	reader := strings.NewReader(testCSV)
	log := logger.New()

	parser := New(log, &readCloser{Reader: reader})

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

type readCloser struct {
	io.Reader
}

func (rc *readCloser) Close() error {
	return nil
}
