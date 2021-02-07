package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/itimofeev/price-store-test/cmd/internal/model"
)

func New(log *logger.Logger, reader io.ReadCloser) *Parser {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'
	return &Parser{
		log:    log,
		reader: csvReader,
		closer: reader,
	}
}

type Parser struct {
	log    *logger.Logger
	closer io.ReadCloser
	reader *csv.Reader
}

// Next returns next parsed struct {product_name;price} from stream
func (p *Parser) Next() (model.ParsedProduct, error) {
	row, err := p.reader.Read()
	if err != nil {
		if closeErr := p.closer.Close(); closeErr != nil {
			p.log.WithError(err).Error("fail to close closer")
		}
		return model.ParsedProduct{}, err
	}

	if len(row) != 2 {
		return model.ParsedProduct{}, errors.New("invalid CSV format, len of columns != 2")
	}

	return model.ParsedProduct{
		Name:  strings.TrimSpace(row[0]),
		Price: strings.TrimSpace(row[1]),
	}, nil
}