package badger

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/apito-io/databasedriver/utility"
	"github.com/dgraph-io/badger/v3"
)

const (
	ProjectCollection = "projects"
	UsersCollection   = "users"
)

func (b *SystemBadgerDriver) SearchResource(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[any], error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {
	return getValue[protobuff.Project](b.Db, ProjectCollection, id)
}

func (b *SystemBadgerDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	return getValue[protobuff.SystemUser](b.Db, UsersCollection, id)
}

func (b *SystemBadgerDriver) GetSystemUserByUsername(ctx context.Context, username string) (*protobuff.SystemUser, error) {
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

func (b *SystemBadgerDriver) CheckProjectName(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) ListProjects(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.Project], error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) ListAllProjects(ctx context.Context, userId string) ([]*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) ListAllUsers(ctx context.Context) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) ListTeams(ctx context.Context, projectId string) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
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

func (b *SystemBadgerDriver) SearchUsers(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.SystemUser], error) {

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

func (b *SystemBadgerDriver) AddSystemUserMetaInfo(ctx context.Context, doc *shared.DefaultDocumentStructure) (*shared.DefaultDocumentStructure, error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) AddTeamMetaInfo(ctx context.Context, docs []*protobuff.SystemUser) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) AddATeamMemberToProject(ctx context.Context, projectId string, memberData map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) RemoveATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) GetATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) UpdateSystemUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	_, err := b.GetSystemUser(ctx, user.Id)
	if err != nil {
		return err
	}
	user.UpdatedAt = utility.GetCurrentTime()
	return b.setValue(UsersCollection, user.Id, user)
}

func (b *SystemBadgerDriver) UpdateProject(ctx context.Context, project *protobuff.Project, replace bool) error {
	_, err := b.GetProject(ctx, project.Id)
	if err != nil {
		return err
	}
	project.UpdatedAt = utility.GetCurrentTime()
	return b.setValue(ProjectCollection, project.Id, project)
}

func (b *SystemBadgerDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) DeleteProjectFromSystem(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) GetSystemUsers(ctx context.Context, keys []string) (map[string]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (b *SystemBadgerDriver) SaveRawData(ctx context.Context, collection string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}
