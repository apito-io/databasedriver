package sql

import (
	"context"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
)

func (p *PostgresSQLDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) GetSystemUserByEmail(ctx context.Context, email string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) CheckProjectName(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) ListProjects(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.Project], error) {
	//TODO implement me
	panic("implement me")
}

/*func (p *PostgresSQLDriver) GetProjectWithRolesAndPermission(ctx context.Context, userId string) ([]*protobuff.ProjectWithRoles, error) {
	//TODO implement me
	panic("implement me")
}*/

func (p *PostgresSQLDriver) ListAllProjects(ctx context.Context, userId string) ([]*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) ListAllUsers(ctx context.Context) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) ListTeams(ctx context.Context, projectId string) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) DeleteWebhook(ctx context.Context, projectId, hookId string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) SearchUsers(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.SystemUser], error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) AddSystemUserMetaInfo(ctx context.Context, doc *shared.DefaultDocumentStructure) (*shared.DefaultDocumentStructure, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) AddTeamMetaInfo(ctx context.Context, docs []*protobuff.SystemUser) ([]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) AddATeamMemberToProject(ctx context.Context, projectId string, memberData map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) RemoveATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) GetATeamMemberFromProject(ctx context.Context, projectId string, memberID string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) CreateProject(ctx context.Context, project *protobuff.Project) (*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) CreateSystemUser(ctx context.Context, user *protobuff.SystemUser) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) UpdateSystemUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) UpdateProject(ctx context.Context, project *protobuff.Project, replace bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) DeleteProjectFromSystem(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) GetSystemUsers(ctx context.Context, keys []string) (map[string]*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresSQLDriver) SaveRawData(ctx context.Context, collection string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}
