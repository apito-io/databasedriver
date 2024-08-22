package boltdb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"go.etcd.io/bbolt"
)

const (
	UserBucket    = "users"
	ProjectBucket = "projects"
)

func (b *SystemBoltDBDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	var user protobuff.SystemUser

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(UserBucket))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		data := bucket.Get([]byte(id))
		if data == nil {
			return errors.New("user not found")
		}
		return json.Unmarshal(data, &user)
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (b *SystemBoltDBDriver) GetSystemUserByUsername(ctx context.Context, username string) (*protobuff.SystemUser, error) {
	var user protobuff.SystemUser

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(UserBucket))
		if bucket == nil {
			return errors.New("bucket not found")
		}

		// Iterate over all users to find the matching username
		return bucket.ForEach(func(k, v []byte) error {
			var u protobuff.SystemUser
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			if u.Username == username {
				user = u
				return nil
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (b *SystemBoltDBDriver) UpdateSystemUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	return b.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(UserBucket))
		if err != nil {
			return err
		}

		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(user.Id), data)
	})
}

func (b *SystemBoltDBDriver) SearchResource(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[any], error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBoltDBDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {
	var project protobuff.Project

	err := b.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(ProjectBucket))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		data := bucket.Get([]byte(id))
		if data == nil {
			return errors.New("project not found")
		}
		return json.Unmarshal(data, &project)
	})

	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (b *SystemBoltDBDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBoltDBDriver) UpdateProject(ctx context.Context, project *protobuff.Project, replace bool) error {
	return b.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(ProjectBucket))
		if err != nil {
			return err
		}

		data, err := json.Marshal(project)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(project.Id), data)
	})
}

func (b *SystemBoltDBDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBoltDBDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBoltDBDriver) DeleteProjectFromSystem(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBoltDBDriver) SaveRawData(ctx context.Context, collection string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}
