package boltdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

func (b *KVBoltDBDriver) AddToSortedSets(ctx context.Context, setName string, key string, exp time.Duration) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(setName))
		if err != nil {
			return err
		}
		expiryTime := time.Now().Add(exp).Unix()
		return bucket.Put([]byte(key), []byte(fmt.Sprintf("%d", expiryTime)))
	})
}

func (b *KVBoltDBDriver) GetFromSortedSets(ctx context.Context, setName string, key string) (float64, error) {
	var val float64
	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(setName))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		data := bucket.Get([]byte(key))
		if data == nil {
			return errors.New("key not found")
		}
		//val = float64(time.Unix(0, 0).Add(time.Duration(string(data))).Unix())
		return nil
	})
	return val, err
}

func (b *KVBoltDBDriver) SetToHashMap(ctx context.Context, hash, key string, value string) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(hash))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (b *KVBoltDBDriver) GetFromHashMap(ctx context.Context, hash, key string) (string, error) {
	var val string
	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hash))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		data := bucket.Get([]byte(key))
		if data == nil {
			return errors.New("key not found")
		}
		val = string(data)
		return nil
	})
	return val, err
}

func (b *KVBoltDBDriver) CheckKeyHashMap(ctx context.Context, hash, key string) bool {
	var exists bool
	b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(hash))
		if bucket != nil && bucket.Get([]byte(key)) != nil {
			exists = true
		}
		return nil
	})
	return exists
}

func (b *KVBoltDBDriver) DelValue(ctx context.Context, key string) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(key))
	})
}

func (b *KVBoltDBDriver) SetValue(ctx context.Context, key string, value string, expiration time.Duration) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (b *KVBoltDBDriver) SetJSONObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	_byte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return b.SetValue(ctx, key, string(_byte), expiration)
}

func (b *KVBoltDBDriver) GetJSONObject(ctx context.Context, key string) (interface{}, error) {
	var jsonData interface{}
	data, err := b.GetValue(ctx, key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (b *KVBoltDBDriver) CheckRedisKey(ctx context.Context, keys ...string) (bool, error) {
	var key string
	switch len(keys) {
	case 2:
		key = fmt.Sprintf("%s:%s", keys[0], keys[1])
	default:
		key = keys[0]
	}
	_, err := b.GetValue(ctx, key)
	if err != nil {
		if errors.Is(err, bolt.ErrBucketNotFound) {
			return false, errors.New("nothing found")
		}
		return false, err
	}
	return true, nil
}

func (b *KVBoltDBDriver) GetValue(ctx context.Context, key string) (string, error) {
	var val string
	err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(key))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		data := bucket.Get([]byte(key))
		if data == nil {
			return errors.New("key not found")
		}
		val = string(data)
		return nil
	})
	return val, err
}

func (b *KVBoltDBDriver) AddToSets(ctx context.Context, key string, value string) error {
	return b.SetToHashMap(ctx, key, value, "")
}

func (b *KVBoltDBDriver) RemoveSets(ctx context.Context, key string, value string) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(key))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		return bucket.Delete([]byte(value))
	})
}

func (b *KVBoltDBDriver) GetStoreDomains(ctx context.Context, sets string, member string) (bool, error) {
	return b.CheckKeyHashMap(ctx, sets, member), nil
}
