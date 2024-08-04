package badger

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/apito-io/databasedriver/utility"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

const (
	ProjectCollection = "projects"
	UsersCollection   = "users"
)

func (b *BadgerDriver) GetOrganizations(ctx context.Context, userId string) (*shared.SearchResponse[protobuff.Organization], error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {
	return getValue[protobuff.Project](b.Db, ProjectCollection, id)
}

func (b *BadgerDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	return getValue[protobuff.SystemUser](b.Db, UsersCollection, id)
}

func (b *BadgerDriver) GetSystemUserByUsername(ctx context.Context, username string) (*protobuff.SystemUser, error) {
	var user *protobuff.SystemUser
	err := b.Db.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 30 // for user search its perfect

		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(UsersCollection)
		var _user protobuff.SystemUser
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &_user)
			})
			if err != nil {
				return err
			}
			if _user.Username == username {
				user = &_user
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (b *BadgerDriver) CheckProjectName(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) ListProjects(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.Project], error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) ListAllProjects(ctx context.Context, userId string) ([]*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) ListAllUsers(ctx context.Context) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) ListTeams(ctx context.Context, projectId string) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
	doc, err := b.GetProject(ctx, param.ProjectId)
	if err != nil {
		return nil, err
	}

	if doc.Schema == nil {
		return nil, errors.New("schema is required")
	}

	return &shared.SearchResponse[protobuff.CloudFunction]{
		Results: doc.Schema.Functions,
	}, nil
}

func (b *BadgerDriver) DeleteWebhook(ctx context.Context, projectId, hookId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) SearchUsers(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.SystemUser], error) {

	var users []*protobuff.SystemUser
	err := b.Db.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 30 // for user search its perfect

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			var _user protobuff.SystemUser
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &_user)
			})
			if err != nil {
				return err
			}
			//if _user.Email
		}

		return nil
	})
	return &shared.SearchResponse[protobuff.SystemUser]{
		Results: users,
	}, err
}

func (b *BadgerDriver) AddSystemUserMetaInfo(ctx context.Context, doc *shared.DefaultDocumentStructure) (*shared.DefaultDocumentStructure, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) AddTeamMetaInfo(ctx context.Context, docs []*protobuff.SystemUser) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) AddATeamMemberToProject(ctx context.Context, projectId string, memberData map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) RemoveATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) GetATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) CreateProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	project.Locals = []string{"en"}
	project.CreatedAt = utility.GetCurrentTime()
	project.UpdatedAt = utility.GetCurrentTime()

	// if starts from example project then transfer the content and model
	/*if project.ProjectTemplate != "" {
		return project, b.TransferSchema(ctx, project.ProjectTemplate, project.Id)
	}*/

	return project, b.setValue(ProjectCollection, project.Id, project)
}

func (b *BadgerDriver) CreateSystemUser(ctx context.Context, user *protobuff.SystemUser) (*protobuff.SystemUser, error) {
	user.XKey = uuid.New().String()
	user.Id = user.XKey
	user.CreatedAt = utility.GetCurrentTime()
	user.UpdatedAt = utility.GetCurrentTime()
	return user, b.setValue(UsersCollection, user.Id, user)
}

func (b *BadgerDriver) UpdateSystemUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	_, err := b.GetSystemUser(ctx, user.Id)
	if err != nil {
		return err
	}
	user.UpdatedAt = utility.GetCurrentTime()
	return b.setValue(UsersCollection, user.Id, user)
}

func (b *BadgerDriver) UpdateProject(ctx context.Context, project *protobuff.Project, replace bool) error {
	_, err := b.GetProject(ctx, project.Id)
	if err != nil {
		return err
	}
	project.UpdatedAt = utility.GetCurrentTime()
	return b.setValue(ProjectCollection, project.Id, project)
}

func (b *BadgerDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) DeleteProjectFromSystem(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) GetSystemUsers(ctx context.Context, keys []string) (map[string]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) SaveRawData(ctx context.Context, collection string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *BadgerDriver) setValue(prefix string, key string, value interface{}) error {
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
