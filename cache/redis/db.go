package redisCache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/redis/go-redis/v9"
)

func (b *CacheDriver) Put(ctx context.Context, id string, cache interface{}) error {
	ttl, _ := strconv.Atoi(b.cfg.CacheTTL)
	data, err := json.Marshal(cache)
	if err != nil {
		fmt.Println("error marshalling cache", err.Error())
		return err
	}
	err = b.Db.Set(ctx, id, data, time.Duration(ttl)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (b *CacheDriver) Get(ctx context.Context, id string) (interface{}, error) {
	val, err := b.Db.Get(ctx, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (b *CacheDriver) ListKeys(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (b *CacheDriver) GetAppCache(ctx context.Context, projectId string) (*shared.ApplicationCache, error) {
	_id := b.idMaker(projectId)

	res, err := b.Db.Get(ctx, _id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var cache *shared.ApplicationCache
	err = json.Unmarshal([]byte(res), &cache)
	if err != nil {
		return nil, err
	}

	// restore the cache
	_val, ok := b.ProjectCache.Load(projectId)
	if ok && cache != nil && _val != nil {
		cache.Project = _val.(*protobuff.Project)
	}

	return cache, err
}

func (b *CacheDriver) PutAppCache(ctx context.Context, projectId string, cache *shared.ApplicationCache) error {
	_id := b.idMaker(projectId)

	b.ProjectCache.Store(projectId, cache.Project)
	_cache := *cache
	_cache.Project = nil

	ttl, _ := strconv.Atoi(b.cfg.CacheTTL)
	_, err := b.Db.Set(ctx, _id, _cache, time.Duration(ttl)*time.Minute).Result()
	if err != nil {
		return err
	}
	return err
}

func (b *CacheDriver) Expire(ctx context.Context, projectId string) error {
	err := b.Db.Del(ctx, projectId).Err()
	if err != nil {
		return err
	}
	// expire the local map
	return nil
}

func (b *CacheDriver) GetProject(ctx context.Context, projectId string) (*protobuff.Project, error) {
	res, err := b.Db.Get(ctx, projectId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var cache *protobuff.Project
	err = json.Unmarshal([]byte(res), &cache)
	if err != nil {
		return nil, err
	}

	return cache, err
}

func (b *CacheDriver) SaveProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	ttl, _ := strconv.Atoi(b.cfg.CacheTTL)
	_, err = b.Db.Set(ctx, project.Id, data, time.Duration(ttl)*time.Second).Result()
	if err != nil {
		return nil, err
	}
	return project, err
}

func (b *CacheDriver) idMaker(projectId string) string {
	return projectId
}
