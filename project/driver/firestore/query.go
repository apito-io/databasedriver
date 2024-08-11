package firestore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/graph-gophers/dataloader"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/tailor-inc/graphql"
	"google.golang.org/api/iterator"
)

type FireStoreDriver struct {
	Db *firestore.Client
}

func (f *FireStoreDriver) RunMigration(ctx context.Context, projectId string) error {
	//TODO implement me
	panic("implement me")
}

func (f *FireStoreDriver) DuplicateModel(ctx context.Context, project *protobuff.Project, modelName, newName string) (*protobuff.ProjectSchema, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FireStoreDriver) GetProjectUsers(ctx context.Context, projectId string, keys []string) (map[string]*shared.DefaultDocumentStructure, error) {
	panic("implement me")
}

func (f *FireStoreDriver) DeleteMediaFile(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) CheckCollectionExists(ctx context.Context, projectId string) (bool, error) {
	// one firebase on project so project collection check is not necessary
	return true, nil
}

func (f *FireStoreDriver) TransferProject(ctx context.Context, userId, from, to string) error {
	return nil
}

func (f *FireStoreDriver) AddCollection(ctx context.Context, projectName string) (*string, error) {
	val, err := f.Db.Collection(projectName).Limit(1).Snapshots(ctx).Next()
	if err != nil {
		return nil, err
	}
	if val.Size > 0 {
		return nil, errors.New("collection Already Exists")
	}
	return &projectName, nil
}

func (f *FireStoreDriver) AddModel(ctx context.Context, project *protobuff.Project, name string, singleRecord bool) (*protobuff.ProjectSchema, error) {
	modelType := &protobuff.ModelType{
		Name: name,
	}

	// if schema not found then create
	if project.Schema == nil {
		project.Schema = &protobuff.ProjectSchema{
			Models: []*protobuff.ModelType{modelType},
		}
	} else {
		var found bool
		for _, ct := range project.Schema.Models {
			if ct.Name == name {
				found = true
				break
			}
		}

		if !found {
			project.Schema.Models = append(project.Schema.Models, modelType)
		} else {
			return nil, errors.New("model Already Defined")
		}
	}

	// check in db label also
	val, err := f.Db.Collection(name).Limit(1).Snapshots(ctx).Next()
	if err != nil {
		return nil, err
	}
	if val.Size > 0 {
		return nil, fmt.Errorf("a model with name `%s` Already Exists in Firebase", name)
	}

	return project.Schema, nil
}

func (f *FireStoreDriver) AddFieldToModel(ctx context.Context, param shared.CommonSystemParams, isUpdate bool, repeatedGroupIdentifier *string) (*protobuff.ModelType, error) {
	if repeatedGroupIdentifier == nil && isUpdate {
		param.Model.Fields = append(param.Model.Fields, param.FieldInfo)
	} else if repeatedGroupIdentifier != nil {
		for _, f := range param.Model.Fields {
			if f.Identifier == *repeatedGroupIdentifier {
				subField := param.FieldInfo
				var found bool
				for i, s := range f.SubFieldInfo {
					if s.Identifier == param.FieldInfo.Identifier {
						f.SubFieldInfo[i] = subField
						found = true
						break
					}
				}
				if !found {
					subField.Serial = uint32(len(f.SubFieldInfo)) + 1
					f.SubFieldInfo = append(f.SubFieldInfo, subField)
				}
			}
		}
	}
	return param.Model, nil
}

func (f *FireStoreDriver) AddRelationFields(ctx context.Context, from *protobuff.ConnectionType, to *protobuff.ConnectionType) error {
	panic("implement me")
}

func (f *FireStoreDriver) ConnectBuilder(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) DisconnectBuilder(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) AddAuthAddOns(ctx context.Context, project *protobuff.Project, auth map[string]interface{}) error {
	panic("implement me")
}

func (f *FireStoreDriver) GetProjectUser(ctx context.Context, phone, email, projectId string) (*shared.DefaultDocumentStructure, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetLoggedInProjectUser(ctx context.Context, param *shared.CommonSystemParams) (*shared.DefaultDocumentStructure, error) {
	panic("implement me")
}

func (f *FireStoreDriver) DeleteDocumentRelation(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) DeleteDocumentsFromProject(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) DropField(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) DropConnections(ctx context.Context, projectId string, from *protobuff.ConnectionType, to *protobuff.ConnectionType) error {
	panic("implement me")
}

func (f *FireStoreDriver) RenameModel(ctx context.Context, project *protobuff.Project, modelName, newName string) error {
	panic("implement me")
}

func (f *FireStoreDriver) ConvertModel(ctx context.Context, project *protobuff.Project, modelName string) error {
	panic("implement me")
}

func (f *FireStoreDriver) RenameField(ctx context.Context, oldFiledName string, repeatedGroup *string, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) GetSystemUser(ctx context.Context, id string) (*protobuff.SystemUser, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetProject(ctx context.Context, id string) (*protobuff.Project, error) {

	var project protobuff.Project
	iter, err := f.Db.Collection("projects").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	err = iter.DataTo(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (f *FireStoreDriver) ListProjects(ctx context.Context, userId string) ([]*protobuff.Project, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetSingleProjectDocumentBytes(ctx context.Context, param shared.CommonSystemParams) ([]byte, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetSingleProjectDocument(ctx context.Context, param shared.CommonSystemParams) (*shared.DefaultDocumentStructure, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetSingleProjectDocumentRevisions(ctx context.Context, param shared.CommonSystemParams) ([]*shared.DocumentRevisionHistory, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetSingleRawDocumentFromProject(ctx context.Context, param shared.CommonSystemParams) (interface{}, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetAllRelationDocumentsOfSingleDocument(ctx context.Context, from string, arg *shared.CommonSystemParams) (interface{}, error) {
	panic("implement me")
}

func (f *FireStoreDriver) ListFunctions(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.CloudFunction], error) {
	panic("implement me")
}

func (f *FireStoreDriver) SearchUsers(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[protobuff.SystemUser], error) {
	panic("implement me")
}

func (f *FireStoreDriver) ListMedias(ctx context.Context, projectId string, param *graphql.ResolveParams) ([]*protobuff.FileDetails, error) {
	panic("implement me")
}

func (f *FireStoreDriver) CountMedias(ctx context.Context, projectId string, param *graphql.ResolveParams) (int, error) {
	panic("implement me")
}

func (f *FireStoreDriver) CountMultiDocumentOfProject(ctx context.Context, param shared.CommonSystemParams, previewMode bool) (int, error) {
	panic("implement me")
}

func (f *FireStoreDriver) QueryMultiDocumentOfProjectBytes(ctx context.Context, param shared.CommonSystemParams) ([]byte, error) {

	var multilineFields []string
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			multilineFields = append(multilineFields, f.Identifier)
		}
	}
	query, err := RootResolverQueryBuilder(param, false)
	if err != nil {
		return nil, err
	}
	collection := f.Db.Collection(param.Model.Name).Query
	for _, q := range query {
		collection = q
	}

	iter := collection.Documents(ctx)

	if err != nil {
		return nil, err
	}
	var docs []*shared.DefaultDocumentStructure
	for {
		rdoc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		var doc shared.DefaultDocumentStructure
		rdoc.DataTo(&doc)

		for _, m := range multilineFields { // #todo if not requestd then dont run
			converter := md.NewConverter("", true, nil)
			if d, ok := doc.Data[m].(map[string]interface{}); ok {
				if html, ok := d["html"].(string); ok {
					markdown, err := converter.ConvertString(html)
					if err != nil {
						fmt.Println(err.Error())
					}
					d["markdown"] = markdown
					d["text"] = strip.StripTags(html)
				}
			}
		}
		docs = append(docs, &doc)
	}

	return []byte{}, nil
}

func (f *FireStoreDriver) QueryMultiDocumentOfProject(ctx context.Context, param shared.CommonSystemParams) ([]*shared.DefaultDocumentStructure, error) {

	var multilineFields []string
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			multilineFields = append(multilineFields, f.Identifier)
		}
	}
	query, err := RootResolverQueryBuilder(param, false)
	if err != nil {
		return nil, err
	}
	collection := f.Db.Collection(param.Model.Name).Query
	for _, q := range query {
		collection = q
	}

	iter := collection.Documents(ctx)

	if err != nil {
		return nil, err
	}
	var docs []*shared.DefaultDocumentStructure
	for {
		rdoc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		var doc shared.DefaultDocumentStructure
		rdoc.DataTo(&doc)

		for _, m := range multilineFields { // #todo if not requestd then dont run
			converter := md.NewConverter("", true, nil)
			if d, ok := doc.Data[m].(map[string]interface{}); ok {
				if html, ok := d["html"].(string); ok {
					markdown, err := converter.ConvertString(html)
					if err != nil {
						fmt.Println(err.Error())
					}
					d["markdown"] = markdown
					d["text"] = strip.StripTags(html)
				}
			}
		}
		docs = append(docs, &doc)
	}

	return docs, nil
}

func (f *FireStoreDriver) AddTeamMetaInfo(ctx context.Context, docs []*protobuff.SystemUser) ([]*protobuff.SystemUser, error) {
	panic("implement me")
}

func (f *FireStoreDriver) AddATeamMemberToProject(ctx context.Context, projectId string, memberData map[string]interface{}) error {
	panic("implement me")
}

func (f *FireStoreDriver) RemoveATeamMemberFromProject(ctx context.Context, projectId string, memberId string) error {
	panic("implement me")
}

func (f *FireStoreDriver) CreateMediaDocument(ctx context.Context, projectId string, media *protobuff.FileDetails) (*protobuff.FileDetails, error) {
	panic("implement me")
}

func (f *FireStoreDriver) UpdateUser(ctx context.Context, user *protobuff.SystemUser, replace bool) error {
	panic("implement me")
}

func (f *FireStoreDriver) CheckTokenBlacklisted(ctx context.Context, tokenId string) error {
	panic("implement me")
}

func (f *FireStoreDriver) BlacklistAToken(ctx context.Context, token map[string]interface{}) error {
	panic("implement me")
}

func (f *FireStoreDriver) UpdateDocumentOfProject(ctx context.Context, param shared.CommonSystemParams, doc *shared.DefaultDocumentStructure, replace bool) error {
	panic("implement me")
}

func (f *FireStoreDriver) DeleteDocumentFromProject(ctx context.Context, param shared.CommonSystemParams) error {
	panic("implement me")
}

func (f *FireStoreDriver) DeleteProject(ctx context.Context, projectId string) error {
	panic("implement me")
}

func (f *FireStoreDriver) CreateRelation(ctx context.Context, projectId string, relation *shared.EdgeRelation) error {
	panic("implement me")
}

func (f *FireStoreDriver) DeleteRelation(ctx context.Context, param *shared.ConnectDisconnectParam, id string) error {
	panic("implement me")
}

func (f *FireStoreDriver) NewInsertableRelations(ctx context.Context, param *shared.ConnectDisconnectParam) ([]string, error) {
	panic("implement me")
}

func (f *FireStoreDriver) CheckOneToOneRelationExists(ctx context.Context, param *shared.ConnectDisconnectParam) (bool, error) {
	panic("implement me")
}

func (f *FireStoreDriver) GetRelationIds(ctx context.Context, param *shared.ConnectDisconnectParam) ([]string, error) {
	panic("implement me")
}

func (f *FireStoreDriver) RelationshipDataLoaderBytes(ctx context.Context, param *shared.CommonSystemParams, connection map[string]interface{}) ([]byte, error) {
	panic("implement me")
}
func (f *FireStoreDriver) RelationshipDataLoader(ctx context.Context, param *shared.CommonSystemParams, connection map[string]interface{}) (interface{}, error) {
	panic("implement me")
}

func (f *FireStoreDriver) MetaDataLoader(ctx context.Context, projectId string, keys *dataloader.Keys) ([]*dataloader.Result, error) {
	panic("implement me")
}

func (f *FireStoreDriver) CountDocOfProjectBytes(ctx context.Context, param *shared.CommonSystemParams) ([]byte, error) {
	panic("implement me")
}

func (f *FireStoreDriver) CountDocOfProject(ctx context.Context, param *shared.CommonSystemParams) (interface{}, error) {
	panic("implement me")
}

func (f *FireStoreDriver) UpdateUsages(ctx context.Context, projectId string, bandwidth float64) error {
	panic("implement me")
}
