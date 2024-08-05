package badger

import (
	badger "github.com/dgraph-io/badger/v3"
	"log"
	"sync"
)

type CacheDriver struct {
	Db           *badger.DB
	ProjectCache sync.Map
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
