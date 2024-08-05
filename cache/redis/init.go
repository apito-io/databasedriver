package redisCache

import (
	"context"
	"fmt"
	"github.com/apito-io/buffers/shared"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
)

type CacheDriver struct {
	cfg *shared.CacheDBConfig
	Db  *redis.Client
	//rh           *rejson.Handler
	ProjectCache sync.Map
}

func GetCacheDriver(cfg *shared.CacheDBConfig) (*CacheDriver, error) {

	dbNo, err := strconv.Atoi(cfg.DB.Database)
	if err != nil {
		return nil, err
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.DB.Host, cfg.DB.Port),
		Password: cfg.DB.Password, // no password set
		DB:       dbNo,            // use default DB
	})
	if err = cli.Ping(context.TODO()).Err(); err != nil {
		return nil, fmt.Errorf(`redis Cache driver error : %s`, err)
	}
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(context.Background(), cli)
	//defer db.Close()

	return &CacheDriver{
		cfg: cfg,
		Db:  cli,
		//rh: rh,
		ProjectCache: sync.Map{},
	}, nil
}
