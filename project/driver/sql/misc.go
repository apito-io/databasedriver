package sql

import "context"

func (S *ProjectSqlDriver) CreateTableOrCollection(ctx context.Context, name string, properties map[string]string) error {
	//TODO implement me
	panic("implement me")
}

func (S *ProjectSqlDriver) DropTableOrCollection(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (S *ProjectSqlDriver) AddDataToTableOrCollection(ctx context.Context, table string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (S *ProjectSqlDriver) UpdateDataToTableOrCollection(ctx context.Context, table string, data map[string]interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (S *ProjectSqlDriver) DeleteDataFromTableOrCollection(ctx context.Context, table string, id string) error {
	//TODO implement me
	panic("implement me")
}
