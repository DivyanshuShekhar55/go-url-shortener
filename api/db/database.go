package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Db_ctx = context.Background()

type ReadDB struct {
	Client *redis.Client
}

func CreateWriteClient(dbNo int) *redis.Client {
	redis_db := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("WRITE_DB_ADDR"),
		Password: os.Getenv("WRITE_DB_PASS"),
		DB:       dbNo,
	})

	return redis_db
}

type WriteDB interface {
	InsertURL(req URL_req, ctx context.Context) (int, error)
	InsertAnalytics(req Analytics_req, ctx context.Context) (int, error)
}

func CreateReadClient(dbNo int) *redis.Client {
	redis_db := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("READ_DB_ADDR"),
		Password: os.Getenv("READ_DB_PASS"),
		DB:       dbNo,
	})

	return redis_db
}

type AnalyticsDB interface {
	StartAnalyticsFlush(shutdown chan struct{})
}

func CreateAnalyticsClient(dbNo int) *redis.Client {
	redis_db := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("ANALYTICS_DB_ADDR"),
		Password: os.Getenv("ANALYTICS_DB_PASS"),
		DB:       dbNo,
	})

	return redis_db
}
