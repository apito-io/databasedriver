package firestore

import "context"

func (f *ProjectFireStoreDriver) CreateTableOrCollection(ctx context.Context, name string, properties map[string]string) error {
	//TODO implement me
	panic("implement me")
}

func (f *ProjectFireStoreDriver) DropTableOrCollection(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}

func (f *ProjectFireStoreDriver) AddDataToTableOrCollection(ctx context.Context, table string, data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (f *ProjectFireStoreDriver) UpdateDataToTableOrCollection(ctx context.Context, table string, data map[string]interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (f *ProjectFireStoreDriver) DeleteDataFromTableOrCollection(ctx context.Context, table string, id string) error {
	//TODO implement me
	panic("implement me")
}
