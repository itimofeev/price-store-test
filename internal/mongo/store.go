// nolint:govet
package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/strategy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/itimofeev/price-store-test/internal/model"
)

func New(url string) *Store {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	breaker, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	if err := retry.Do(breaker, func(ctx context.Context) error {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		return client.Ping(ctx, readpref.Primary())
	}, strategy.Limit(5), strategy.Backoff(backoff.Linear(time.Second))); err != nil {
		panic(err)
	}

	return &Store{
		client: client,
	}
}

type Store struct {
	client *mongo.Client
}

func (s *Store) SaveProduct(ctx context.Context, updateTime time.Time, product model.ParsedProduct) (saved model.Product, err error) {
	filter := bson.D{{"name", product.Name}}
	opts := options.Update().SetUpsert(true)

	update := bson.D{
		{"$set", bson.D{{"price", product.Price}, {"lastUpdate", updateTime}}},
		{"$inc", bson.D{{"updateCount", 1}}},
	}

	_, err = s.productsCollection().UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return saved, fmt.Errorf("failed to save product in mongo: %w", err)
	}

	cur, err := s.productsCollection().Find(ctx, filter)
	if err != nil {
		return saved, fmt.Errorf("failed to find product by name: %w", err)
	}
	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		return saved, fmt.Errorf("product not found by name: %w", err)
	}
	if err := cur.Err(); err != nil {
		return saved, fmt.Errorf("failed cursor error: %w", err)
	}

	return decodeProduct(cur)
}

func (s *Store) ListProducts(ctx context.Context, order string, limit, offset int) (products []model.Product, err error) {
	limit64 := int64(limit)
	offset64 := int64(offset)
	opts := &options.FindOptions{
		Limit: &limit64,
		Skip:  &offset64,
		Sort:  bson.D{{"lastUpdate", -1}},
	}
	cur, err := s.productsCollection().Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		saved, err := decodeProduct(cur)
		if err != nil {
			return nil, err
		}
		products = append(products, saved)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("failed cursor error: %w", err)
	}
	return products, nil
}

func (s *Store) productsCollection() *mongo.Collection {
	return s.client.Database("db").Collection("products")
}

func decodeProduct(cur *mongo.Cursor) (product model.Product, err error) {
	var result bson.D
	if err := cur.Decode(&result); err != nil {
		return product, fmt.Errorf("failed to decode result from mongo: %w", err)
	}
	m := result.Map()
	product.ID = m["_id"].(primitive.ObjectID).String()
	product.Price = decodeInt(m["price"])
	product.Name = m["name"].(string)
	lastUpdateDT := m["lastUpdate"].(primitive.DateTime)
	product.LastUpdate = lastUpdateDT.Time()
	product.UpdateCount = decodeInt(m["updateCount"]) - 1

	return product, nil
}

func decodeInt(i interface{}) int64 {
	if d, ok := i.(int64); ok {
		return d
	}
	if d, ok := i.(int32); ok {
		return int64(d)
	}
	if d, ok := i.(int); ok {
		return int64(d)
	}
	panic("unknown type of int")
}
