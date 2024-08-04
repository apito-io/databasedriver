package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/apito-io/buffers/shared"
	"github.com/redis/go-redis/v9"
)

// #todo redis sentinal service

type KVRedisService struct {
	client *redis.Client
}

func GetKVRedisDriver(ctx context.Context, db *shared.CommonDatabaseConfig) (*KVRedisService, error) {

	dbNo, err := strconv.Atoi(db.Database)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", db.Host, db.Port),
		Password: db.Password, // no password set
		DB:       dbNo,        // use default DB
	})

	err = client.Ping(ctx).Err()
	if err != nil {
		return nil, errors.New(fmt.Sprintf(`redis KV driver error : %s`, err))
	}

	return &KVRedisService{
		client: client,
	}, nil
}
