package badger

import (
	"context"
	"encoding/json"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/databasedriver/utility"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

func (b *SystemBadgerDriver) createProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	project.Locals = []string{"en"}
	project.CreatedAt = utility.GetCurrentTime()
	project.UpdatedAt = utility.GetCurrentTime()

	// if starts from example project then transfer the content and model
	/*if project.ProjectTemplate != "" {
		return project, b.TransferSchema(ctx, project.ProjectTemplate, project.Id)
	}*/

	return project, b.setValue(ProjectCollection, project.ID, project)
}

func (b *SystemBadgerDriver) createSystemUser(ctx context.Context, user *protobuff.SystemUser) (*protobuff.SystemUser, error) {
	user.XKey = uuid.New().String()
	user.ID = user.XKey
	user.CreatedAt = utility.GetCurrentTime()
	user.UpdatedAt = utility.GetCurrentTime()
	return user, b.setValue(UsersCollection, user.ID, user)
}

func (b *SystemBadgerDriver) setValue(prefix string, key string, value interface{}) error {
	return b.Db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(value)
		if err != nil {
			return nil
		}
		e := badger.NewEntry([]byte(prefix+"_"+key), data)
		return txn.SetEntry(e)
	})
}

func getValue[T any](db *badger.DB, prefix string, key string) (*T, error) {
	var res T
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(prefix + "_" + key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &res)
		})
		if err != nil {
			return err
		}
		return nil
	})
	return &res, err
}

func filter[T any](filter map[string]interface{}) (*T, error) {
	var data T

	// using golang reflect iterate through all the fields of T and match
	//T.field = map key, value

	return &data, nil
}
