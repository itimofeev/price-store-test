package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/itimofeev/price-store-test/internal/model"
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

const numberOfColumns = 2

// regexp allows 0-2 signs after delimiter
// allowed: .0, 0.10
// not allowed: .
var floatRegexp = regexp.MustCompile(`^\d*(\.\d{1,2})?$`)

// Next returns next parsed struct {product_name;price} from stream
func (p *Parser) Next() (model.ParsedProduct, error) {
	row, err := p.reader.Read()
	if err != nil {
		if closeErr := p.closer.Close(); closeErr != nil {
			p.log.WithError(err).Error("fail to close closer")
		}
		return model.ParsedProduct{}, err
	}

	if len(row) != numberOfColumns {
		return model.ParsedProduct{}, errors.New("invalid CSV format, len of columns != 2")
	}

	priceStr := strings.TrimSpace(row[1])
	if !floatRegexp.MatchString(priceStr) {
		return model.ParsedProduct{}, fmt.Errorf("invalid price: %s", row[1])
	}
	priceStr = strings.Replace(priceStr, ".", "", 1)
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		return model.ParsedProduct{}, fmt.Errorf("unable to parse price: %w", err)
	}

	return model.ParsedProduct{
		Name:  strings.TrimSpace(row[0]),
		Price: price,
	}, nil
}
