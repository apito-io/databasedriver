package firestore

import (
	"context"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
)

func (f *ProjectFireStoreDriver) RemoveAuthAddOns(ctx context.Context, project *protobuff.Project, option map[string]interface{}) error {

	return nil
}

func (f *ProjectFireStoreDriver) AddDocumentToProject(ctx context.Context, projectId string, modelName string, doc *shared.DefaultDocumentStructure) (interface{}, error) {
	_, err := f.Db.Collection(modelName).Doc(doc.Id).Set(ctx, doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
