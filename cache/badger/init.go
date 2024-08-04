package badger

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/goccy/go-json"
)

type CacheDriver struct {
	Db           *badger.DB
	ProjectCache sync.Map
}

func (b *CacheDriver) Put(ctx context.Context, id string, cache interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *CacheDriver) Get(ctx context.Context, id string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func GetCacheDriver() (*CacheDriver, error) {

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("cache/storage"))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	return &CacheDriver{Db: db, ProjectCache: sync.Map{}}, nil
}

func (b *CacheDriver) ListKeys(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (b *CacheDriver) GetAppCache(ctx context.Context, projectId string) (*shared.ApplicationCache, error) {
	_id := b.idMaker(projectId)
	var cache *shared.ApplicationCache
	err := b.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(_id))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &cache)
		})
		if err != nil {
			return err
		}
		return nil
	})

	// restore the cache
	_val, ok := b.ProjectCache.Load(projectId)
	if ok && cache != nil && _val != nil {
		cache.Project = _val.(*protobuff.Project)
	}

	return cache, err
}

func (b *CacheDriver) PutAppCache(ctx context.Context, projectId string, cache *shared.ApplicationCache) error {
	_id := b.idMaker(projectId)
	err := b.Db.Update(func(txn *badger.Txn) error {
		b.ProjectCache.Store(projectId, cache.Project)
		_cache := *cache
		_cache.Project = nil
		data, err := json.Marshal(_cache)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		e := badger.NewEntry([]byte(_id), data).WithTTL(1 * time.Hour)
		err = txn.SetEntry(e)
		return err
	})
	return err
}

func (b *CacheDriver) Expire(ctx context.Context, projectId string) error {
	_id := b.idMaker(projectId)
	err := b.Db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal([]byte{})
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(_id), data).WithTTL(1 * time.Millisecond)
		err = txn.SetEntry(e)
		return err
	})
	// expire the local map
	b.ProjectCache.Delete(projectId)
	return err
}

func (b *CacheDriver) GetProject(ctx context.Context, projectId string) (*protobuff.Project, error) {
	var project *protobuff.Project
	err := b.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(projectId))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &project)
		})
		if err != nil {
			return err
		}
		return nil
	})

	return project, err
}

func (b *CacheDriver) SaveProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	err := b.Db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(project)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(project.Id), data).WithTTL(1 * time.Minute)
		err = txn.SetEntry(e)
		return err
	})
	return project, err
}

func (b *CacheDriver) idMaker(projectId string) string {
	return fmt.Sprintf(`%s`, projectId)
}
