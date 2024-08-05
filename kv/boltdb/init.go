package boltdb

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

type KVBoltDBDriver struct {
	DB *bolt.DB
}

func GetKVBadgerDriver() (*KVBoltDBDriver, error) {

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("./db/kv.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	return &KVBoltDBDriver{DB: db}, nil
}
