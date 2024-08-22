package sql

import (
	"context"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
)

func (p *SystemSQLDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) GetSystemUserByUsername(ctx context.Context, username string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) UpdateSystemUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) SearchResource(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[any], error) {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) UpdateProject(ctx context.Context, project *protobuff.Project, replace bool) error {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) DeleteProjectFromSystem(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (p *SystemSQLDriver) SaveRawData(ctx context.Context, collection string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}
