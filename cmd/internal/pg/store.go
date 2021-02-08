package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	logger "github.com/sirupsen/logrus"

	"github.com/itimofeev/price-store-test/cmd/internal/model"
	"github.com/itimofeev/price-store-test/cmd/internal/util"
)

func New(log *logger.Logger, url string) *Store {
	opts, err := pg.ParseURL(url)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opts)
	db.AddQueryHook(dbLogger{log: log})
	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	if err := doMigrationIfNeeded(db); err != nil {
		panic(err)
	}
	return &Store{db: db}
}

type Store struct {
	db *pg.DB
}

func (s *Store) SaveProduct(ctx context.Context, updateTime time.Time, product model.ParsedProduct) (saved model.Product, err error) {
	sql := `
		INSERT INTO
			products (id, name, price, last_update)
		VALUES
			(?, ?, ?, ?)
		ON CONFLICT (name) DO UPDATE
			SET
				price        = excluded.price,
				last_update  = excluded.last_update,
				update_count = products.update_count + 1
		RETURNING *
	`

	_, err = s.db.WithContext(ctx).QueryOne(&saved, sql, util.RandomID(), product.Name, product.Price, updateTime)
	if err != nil {
		return saved, fmt.Errorf("failed to save product: %w", err)
	}
	return saved, nil
}

func (s *Store) ListProducts(ctx context.Context, order string, limit, offset int) (products []model.Product, err error) {
	query := s.db.WithContext(ctx).
		Model(&products).
		Order(order).
		Limit(limit).
		Offset(offset)
	if err = query.Select(); err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	return products, nil
}

type dbLogger struct {
	log *logger.Logger
}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	if d.log != nil {
		query, _ := q.FormattedQuery()
		d.log.WithField("query", string(query)).Debug("query log")
	}
	return nil
}
