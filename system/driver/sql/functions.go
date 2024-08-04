package sql

import (
	"context"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
)

func (p PostgreSQLDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) GetSystemUserByEmail(ctx context.Context, email string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) CheckProjectName(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) ListProjects(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.Project], error) {
	//TODO implement me
	panic("implement me")
}

/*func (p PostgreSQLDriver) GetProjectWithRolesAndPermission(ctx context.Context, userId string) ([]*protobuff.ProjectWithRoles, error) {
	//TODO implement me
	panic("implement me")
}*/

func (p PostgreSQLDriver) ListAllProjects(ctx context.Context, userId string) ([]*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) ListAllUsers(ctx context.Context) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) ListTeams(ctx context.Context, projectId string) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) DeleteWebhook(ctx context.Context, projectId, hookId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) SearchUsers(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.SystemUser], error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) AddSystemUserMetaInfo(ctx context.Context, doc *shared.DefaultDocumentStructure) (*shared.DefaultDocumentStructure, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) AddTeamMetaInfo(ctx context.Context, docs []*protobuff.SystemUser) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) AddATeamMemberToProject(ctx context.Context, projectId string, memberData map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) RemoveATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) GetATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) CreateProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) CreateSystemUser(ctx context.Context, user *protobuff.SystemUser) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) UpdateSystemUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) UpdateProject(ctx context.Context, project *protobuff.Project, replace bool) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) DeleteProjectFromSystem(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) GetSystemUsers(ctx context.Context, keys []string) (map[string]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) SaveRawData(ctx context.Context, collection string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}
