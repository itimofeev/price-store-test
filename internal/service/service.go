package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/itimofeev/price-store-test/internal/csv"
	"github.com/itimofeev/price-store-test/internal/model"
)

type Downloader interface {
	GetCSV(ctx context.Context, url string) (io.ReadCloser, error)
}

type Store interface {
	SaveProduct(ctx context.Context, updateTime time.Time, product model.ParsedProduct) (saved model.Product, err error)
	ListProducts(ctx context.Context, order string, limit, offset int) (products []model.Product, err error)
}

func New(log *logrus.Logger, d Downloader, s Store) *Service {
	return &Service{
		log: log,
		d:   d,
		s:   s,
	}
}

type Service struct {
	log *logrus.Logger
	d   Downloader
	s   Store
}

func (s *Service) ProcessCSV(ctx context.Context, url string) error {
	now := time.Now()
	csvStream, err := s.d.GetCSV(ctx, url)
	if err != nil {
		return fmt.Errorf("failed during getting csv: %w", err)
	}

	csvParser := csv.New(s.log, csvStream)

	for {
		parsed, err := csvParser.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed during parsing csv: %w", err)
		}

		if _, err := s.s.SaveProduct(ctx, now, parsed); err != nil {
			return fmt.Errorf("failed to save product: %w", err)
		}
	}

	return nil
}

func (s *Service) ListProducts(ctx context.Context, order string, limit, offset int) ([]model.Product, error) {
	return s.s.ListProducts(ctx, order, limit, offset)
}
