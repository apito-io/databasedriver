package boltdb

import (
	"context"
	"encoding/json"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/databasedriver/utility"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func (b *SystemBoltDBDriver) createProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	project.Locals = []string{"en"}
	project.CreatedAt = utility.GetCurrentTime()
	project.UpdatedAt = utility.GetCurrentTime()

	// if starts from example project then transfer the content and model
	/*if project.ProjectTemplate != "" {
		return project, b.TransferSchema(ctx, project.ProjectTemplate, project.Id)
	}*/

	return project, b.setValue(ProjectBucket, project.ID, project)
}

func (b *SystemBoltDBDriver) createSystemUser(ctx context.Context, user *protobuff.SystemUser) (*protobuff.SystemUser, error) {
	user.XKey = uuid.New().String()
	user.ID = user.XKey
	user.CreatedAt = utility.GetCurrentTime()
	user.UpdatedAt = utility.GetCurrentTime()
	return user, b.setValue(UserBucket, user.ID, user)
}

func getValue[T any](db *bbolt.DB, bucketName string, key string) (*T, error) {
	var result T

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return bbolt.ErrBucketNotFound
		}

		val := bucket.Get([]byte(key))
		if val == nil {
			return bbolt.ErrBucketNotFound
		}

		return json.Unmarshal(val, &result)
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *SystemBoltDBDriver) setValue(bucketName string, key string, value interface{}) error {
	return b.DB.Update(func(tx *bbolt.Tx) error {
		// Create or access the existing bucket
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		// Marshal the value to JSON
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}

		// Store the value in the bucket using the key
		return bucket.Put([]byte(key), data)
	})
}
