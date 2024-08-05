package memoryCache

import (
	"context"
	"errors"
	"fmt"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
)

func (b *CacheDriver) Put(ctx context.Context, id string, cache interface{}) error {
	b.Cache.Store(id, cache)
	return nil
}

func (b *CacheDriver) Get(ctx context.Context, id string) (interface{}, error) {
	if val, ok := b.Cache.Load(id); ok && val != nil {
		return val, nil
	}
	return nil, errors.New("cache not found. fetch one")
}

func (b *CacheDriver) ListKeys(ctx context.Context) ([]string, error) {
	var _keys []string
	b.Cache.Range(func(key, value interface{}) bool {
		_keys = append(_keys, key.(string))
		return true
	})
	return _keys, nil
}

func (b *CacheDriver) GetAppCache(ctx context.Context, projectId string) (*shared.ApplicationCache, error) {
	return nil, nil
}

func (b *CacheDriver) PutAppCache(ctx context.Context, projectId string, cache *shared.ApplicationCache) error {
	return nil
}

func (b *CacheDriver) Expire(ctx context.Context, id string) error {
	b.Cache.Delete(id)
	return nil
}

func (b *CacheDriver) GetProject(ctx context.Context, projectId string) (*protobuff.Project, error) {
	var _project *protobuff.Project
	if val, ok := b.Cache.Load(projectId); ok && val != nil {
		_project = val.(*protobuff.Project)
		return _project, nil
	}
	return nil, errors.New("project cache not found. fetch one")
}

func (b *CacheDriver) SaveProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	b.Cache.Store(project.Id, project)
	return project, nil
}

func (b *CacheDriver) idMaker(projectId string) string {
	return fmt.Sprintf(`%s`, projectId)
}
