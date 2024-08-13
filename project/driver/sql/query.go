package sql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	_const "github.com/apito-io/databasedriver"
	"github.com/jinzhu/inflection"
	_ "github.com/lib/pq"
	"github.com/tailor-inc/graphql"
	"github.com/uptrace/bun"
	"gorm.io/gorm"
)

func (S *SqlDriver) GetProjectUsers(ctx context.Context, projectId string, keys []string) (map[string]*shared.DefaultDocumentStructure, error) {
	panic("get project users not implemented")
}

func (S *SqlDriver) CountDocOfProject(ctx context.Context, param *shared.CommonSystemParams) (interface{}, error) {
	query, err := RootConnectionResolverQueryBuilder(param)
	if err != nil {
		return nil, err
	}

	res, err := S.ORM.ExecContext(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(res.RowsAffected())

	return map[string]interface{}{
		"total": 0,
	}, nil
}

func (S *SqlDriver) CountDocOfProjectBytes(ctx context.Context, param *shared.CommonSystemParams) ([]byte, error) {
	query, err := RootConnectionResolverQueryBuilder(param)
	if err != nil {
		return nil, err
	}

	_, err = S.ORM.Exec(query, nil)
	if err != nil {
		return nil, err
	}

	return []byte{}, nil
}

func (S *SqlDriver) AddAuthAddOns(ctx context.Context, project *protobuff.Project, auth map[string]interface{}) error {
	panic("add auth addons not implemented")
}

func (S *SqlDriver) ConnectBuilder(ctx context.Context, param shared.CommonSystemParams) error {
	var err error
	for _, param := range param.ConDisParam {
		for _, id := range param.ConnectionIds {
			switch param.ConnectionType {
			case "forward":
				tableName := inflection.Plural(param.ForwardConnectionType.Model)
				switch param.BackwardConnectionType.Relation {
				case "has_one":
					model := map[string]interface{}{
						fmt.Sprintf(`%s_id`, param.BackwardConnectionType.Model): param.ForwardConnectionId,
					}
					_, err := S.ORM.NewUpdate().Table(tableName).Where("id = ?", id).Model(model).Exec(ctx)
					if err != nil {
						return err
					}
					break
				case "has_many":
					tableName = fmt.Sprintf(`%s_%s`, inflection.Plural(param.BackwardConnectionType.Model), tableName)
					model := map[string]interface{}{
						fmt.Sprintf(`%s_id`, param.BackwardConnectionType.Model): param.ForwardConnectionId,
						fmt.Sprintf(`%s_id`, param.ForwardConnectionType.Model):  id,
					}
					_, err = S.ORM.NewInsert().Table(tableName).Model(model).Exec(ctx)
					if err != nil {
						return err
					}
					break
				}
				break
			case "backward":
				tableName := inflection.Plural(param.ForwardConnectionType.Model)
				switch param.ForwardConnectionType.Relation {
				case "has_one":
					model := map[string]interface{}{
						fmt.Sprintf(`%s_id`, param.BackwardConnectionType.Model): param.ForwardConnectionId,
					}
					_, err = S.ORM.NewUpdate().Table(tableName).Where("id = ?", id).Model(model).Exec(ctx)
					if err != nil {
						return err
					}
					break
				case "has_many":
					tableName = fmt.Sprintf(`%s_%s`, inflection.Plural(param.BackwardConnectionType.Model), tableName)
					u := map[string]interface{}{
						fmt.Sprintf(`%s_id`, param.BackwardConnectionType.Model): param.ForwardConnectionId,
						fmt.Sprintf(`%s_id`, param.ForwardConnectionType.Model):  id,
					}
					_, err = S.ORM.NewInsert().Table(tableName).Model(u).Exec(ctx)
					if err != nil {
						return err
					}
					break
				}
				break
			}
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (S *SqlDriver) DisconnectBuilder(ctx context.Context, param shared.CommonSystemParams) error {
	return nil
}

func (S *SqlDriver) CheckCollectionExists(ctx context.Context, projectId string) (bool, error) {
	return true, nil
}

func (S *SqlDriver) CheckDBExists(ctx context.Context, projectId string) (bool, error) {
	var err error
	var foundDbName string

	switch S.DriverCredential.Engine {
	case _const.MySQLDriver:
		err = S.ORM.NewSelect().
			Column("SCHEMA_NAME").
			Table("INFORMATION_SCHEMA.SCHEMATA").
			Where("SCHEMA_NAME = ?", projectId).
			Scan(ctx, &foundDbName)
	case _const.PostgresSQLDriver:
		// SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('dbname');
		err = S.ORM.NewSelect().
			Column("datname").
			Table("pg_catalog.pg_database").
			Where("lower(datname) = lower(?)", projectId).
			Scan(ctx, &foundDbName)
	case _const.SQLServerDriver:
		err = S.ORM.NewSelect().
			Column("name").
			Table("sys.databases").
			Where("lower(name) = lower(?)", projectId).
			Scan(ctx, &foundDbName)
	case _const.SQLiteDriver:
		var exists bool
		err = S.ORM.NewSelect().
			ColumnExpr("count(*) > 0").
			Table("sqlite_master").
			Where("type = 'table' AND name = 'some_known_table'").
			Scan(ctx, &exists)
		if exists {
			return true, nil
		}
	}

	if err != nil {
		return false, err
	}
	if foundDbName == projectId {
		return true, nil
	}
	return false, nil
}

func (S *SqlDriver) GetProjectUser(ctx context.Context, phone, email, projectId string) (*shared.DefaultDocumentStructure, error) {
	panic("get project user not implemented")
}

func (S *SqlDriver) GetLoggedInProjectUser(ctx context.Context, param *shared.CommonSystemParams) (*shared.DefaultDocumentStructure, error) {
	panic("get logged in project user")
}

/* deprecated
func (S *SqlDriver) GetAllPreviewDocumentsByModel(param shared.CommonSystemParams) ([]*protobuff.PreviewMode, error) {
	query, err := RootResolverQueryBuilder(param, true)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	err = S.Gorm.Raw(*query).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*protobuff.PreviewMode{}, nil
		} else {
			return nil, err
		}
	}

	var docs []*protobuff.PreviewMode
	for _, res := range results {
		doc := &protobuff.PreviewMode{}
		if val, ok := res["id"].([]byte); ok {
			doc.Id = string(val)
		}
		if val, ok := res["title"].(string); ok {
			doc.Title = val
		}
		if val, ok := res["status"].(string); ok {
			doc.Status = val
		} else {
			doc.Status = "draft" // default
		}
		if val, ok := res["icon"].(string); ok {
			doc.Icon = val
		}
		// filter doc title
		title := strip.StripTags(doc.Title)
		if len(title) > 35 {
			doc.Title = title[0:35] + "..."
		} else {
			doc.Title = title
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
*/

func (S *SqlDriver) GetSingleProjectDocumentRevisions(ctx context.Context, param shared.CommonSystemParams) ([]*shared.DocumentRevisionHistory, error) {
	panic("get single project document revision not implemented")
}

func (S *SqlDriver) GetSingleProjectDocumentBytes(ctx context.Context, param shared.CommonSystemParams) ([]byte, error) {

	var local string
	if val, ok := param.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	returnType := SelectBuilder("y", local, param.Model, false)

	tableName := inflection.Plural(param.Model.Name)
	result := map[string]interface{}{}
	err := S.ORM.NewRaw(`SELECT ? FROM ? AS x LEFT JOIN meta as y on y.doc_id = x.id WHERE x.id = '?'`,
		strings.Join(returnType, ", "),
		bun.Ident(tableName),
		param.DocumentId,
	).Scan(ctx, &result)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []byte{}, nil
		} else {
			return nil, err
		}
	}

	classification := FieldClassification{}
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			classification.MultilineFields = append(classification.MultilineFields, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && f.Validation.IsGallery {
			classification.GalleryField = append(classification.GalleryField, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && !f.Validation.IsGallery {
			classification.PictureField = append(classification.PictureField, f.Identifier)
		} else if f.FieldType == "number" && f.InputType == "double" {
			classification.DoubleFields = append(classification.DoubleFields, f.Identifier)
		} else if f.FieldType == "list" && f.Validation != nil && (len(f.Validation.FixedListElements) == 0 || f.Validation.IsMultiChoice) {
			classification.ListFields = append(classification.ListFields, f.Identifier)
		} else if f.FieldType == "repeated" {
			classification.RepeatedFields = map[string][]*protobuff.FieldInfo{
				f.Identifier: f.SubFieldInfo,
			}
		}
	}

	doc, err := CommonDocTransformation(param.Model, local, result, &classification)
	if err != nil {
		return nil, err
	}

	doc.Type = param.Model.Name
	return nil, nil
}

func (S *SqlDriver) GetSingleProjectDocument(ctx context.Context, param shared.CommonSystemParams) (*shared.DefaultDocumentStructure, error) {

	var local string
	if val, ok := param.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	returnType := SelectBuilder("y", local, param.Model, false)

	tableName := inflection.Plural(param.Model.Name)
	result := map[string]interface{}{}
	err := S.ORM.NewRaw(`SELECT ? FROM ? AS x LEFT JOIN meta as y on y.doc_id = x.id WHERE x.id = '?'`,
		strings.Join(returnType, ", "),
		bun.Ident(tableName),
		param.DocumentId,
	).Scan(ctx, &result)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &shared.DefaultDocumentStructure{}, nil
		} else {
			return nil, err
		}
	}

	classification := FieldClassification{}
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			classification.MultilineFields = append(classification.MultilineFields, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && f.Validation.IsGallery {
			classification.GalleryField = append(classification.GalleryField, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && !f.Validation.IsGallery {
			classification.PictureField = append(classification.PictureField, f.Identifier)
		} else if f.FieldType == "number" && f.InputType == "double" {
			classification.DoubleFields = append(classification.DoubleFields, f.Identifier)
		} else if f.FieldType == "list" && f.Validation != nil && (len(f.Validation.FixedListElements) == 0 || f.Validation.IsMultiChoice) {
			classification.ListFields = append(classification.ListFields, f.Identifier)
		} else if f.FieldType == "repeated" {
			classification.RepeatedFields = map[string][]*protobuff.FieldInfo{
				f.Identifier: f.SubFieldInfo,
			}
		}
	}

	doc, err := CommonDocTransformation(param.Model, local, result, &classification)
	if err != nil {
		return nil, err
	}

	doc.Type = param.Model.Name
	return doc, nil
}

func (S *SqlDriver) GetSingleRawDocumentFromProject(ctx context.Context, param shared.CommonSystemParams) (interface{}, error) {
	if param.SinglePageData {
		return &shared.DefaultDocumentStructure{
			Key:  param.DocumentId,
			Id:   param.DocumentId,
			Type: param.Model.Name,
			Data: map[string]interface{}{},
			Meta: &protobuff.MetaField{
				LastModifiedBy: &protobuff.SystemUser{},
			},
		}, nil
	}

	var local string
	if val, ok := param.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	returnType := fmt.Sprintf(`
	 x.*, y.created_at AS sys_created_at, y.updated_at AS sys_updated_at,
		y.status as sys_status
	`)

	tableName := inflection.Plural(param.Model.Name)
	result := map[string]interface{}{}
	err := S.ORM.NewRaw(`SELECT ? FROM ? AS x LEFT JOIN meta as y on y.doc_id = x.id WHERE x.id = '?'`,
		returnType,
		bun.Ident(tableName),
		param.DocumentId,
	).Scan(ctx, &result)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &shared.DefaultDocumentStructure{}, nil
		} else {
			return nil, err
		}
	}

	classification := FieldClassification{}
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			classification.MultilineFields = append(classification.MultilineFields, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && f.Validation.IsGallery {
			classification.GalleryField = append(classification.GalleryField, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && !f.Validation.IsGallery {
			classification.PictureField = append(classification.PictureField, f.Identifier)
		} else if f.FieldType == "number" && f.InputType == "double" {
			classification.DoubleFields = append(classification.DoubleFields, f.Identifier)
		} else if f.FieldType == "list" && f.Validation != nil && (len(f.Validation.FixedListElements) == 0 || f.Validation.IsMultiChoice) {
			classification.ListFields = append(classification.ListFields, f.Identifier)
		} else if f.FieldType == "repeated" {
			classification.RepeatedFields = map[string][]*protobuff.FieldInfo{
				f.Identifier: f.SubFieldInfo,
			}
		}
	}

	doc, err := CommonDocTransformation(param.Model, local, result, &classification)
	if err != nil {
		return nil, err
	}

	doc.Type = param.Model.Name
	return doc, nil
}

func (S *SqlDriver) GetAllRelationDocumentsOfSingleDocument(ctx context.Context, from string, arg *shared.CommonSystemParams) (interface{}, error) {
	// query relations and find all docs
	arg.DocumentIDs = []string{arg.DocumentId}
	arg.OnlyReturnCount = true
	query, relationType, err := BuildCombinedRelationQuery("", from, arg)
	if err != nil {
		return nil, err
	}

	switch *relationType {
	case "has_many":
		var result []map[string]interface{}
		err = S.ORM.NewRaw(*query).Scan(ctx, &result)
		if err != nil {
			return nil, err
		}
		var docs []*protobuff.PreviewMode
		for _, res := range result {
			doc := protobuff.PreviewMode{}
			if val, ok := res["id"].([]byte); ok {
				id := string(val)
				doc.Id = id
			} else {
				// if no id then return nil
				return []*protobuff.PreviewMode{}, nil
			}
			if val, ok := res["title"].(string); ok {
				doc.Title = val
			}
			if val, ok := res["icon"].(string); ok {
				doc.Icon = val
			}
			if val, ok := res["status"].(string); ok {
				doc.Status = val
			}
			docs = append(docs, &doc)
		}
		return docs, nil
	case "has_one":
		result := map[string]interface{}{}
		err = S.ORM.NewRaw(*query).Scan(ctx, &result)
		if err != nil {
			return nil, err
		}
		doc := protobuff.PreviewMode{}
		if val, ok := result["id"].([]byte); ok {
			id := string(val)
			doc.Id = id
		} else {
			// if no id then return nil
			return nil, nil
		}
		if val, ok := result["title"].(string); ok {
			doc.Title = val
		}
		if val, ok := result["icon"].(string); ok {
			doc.Icon = val
		}
		if val, ok := result["status"].(string); ok {
			doc.Status = val
		}
		return doc, nil
	}

	return nil, errors.New("invalid Relation Type")
}

func (S *SqlDriver) CountMedias(ctx context.Context, projectId string, param *graphql.ResolveParams) (int, error) {
	return 0, nil
}

func (S *SqlDriver) ListMedias(ctx context.Context, projectId string, param *graphql.ResolveParams) ([]*protobuff.FileDetails, error) {
	query, err := RootResolverMediaQueryBuilder(param)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = S.ORM.NewRaw(query).Scan(ctx, &result)
	if err != nil {
		return nil, err
	}

	var docs []*protobuff.FileDetails
	for _, res := range result {
		doc, err := MediaDocTransformation("media", res)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (S *SqlDriver) CountMultiDocumentOfProject(ctx context.Context, param shared.CommonSystemParams, previewMode bool) (int, error) {

	query, err := RootResolverQueryBuilder(param, true)
	if err != nil {
		return 0, err
	}

	var result int64
	err = S.ORM.NewRaw(*query).Scan(ctx, &result)
	if err != nil {
		return 0, err
	}

	return int(result), nil

}

func (S *SqlDriver) QueryMultiDocumentOfProjectBytes(ctx context.Context, param shared.CommonSystemParams) ([]byte, error) {

	var local string
	if val, ok := param.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	query, err := RootResolverQueryBuilder(param, false)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = S.ORM.NewRaw(*query).Scan(ctx, &result)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // send empty response for table
			return []byte{}, nil
		}
		return nil, err
	}

	classification := FieldClassification{}
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			classification.MultilineFields = append(classification.MultilineFields, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && f.Validation.IsGallery {
			classification.GalleryField = append(classification.GalleryField, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && !f.Validation.IsGallery {
			classification.PictureField = append(classification.PictureField, f.Identifier)
		} else if f.FieldType == "number" && f.InputType == "double" {
			classification.DoubleFields = append(classification.DoubleFields, f.Identifier)
		} else if f.FieldType == "list" && f.Validation != nil && (len(f.Validation.FixedListElements) == 0 || f.Validation.IsMultiChoice) {
			classification.ListFields = append(classification.ListFields, f.Identifier)
		}
	}

	var docs []*shared.DefaultDocumentStructure
	for _, res := range result {
		doc, err := CommonDocTransformation(param.Model, local, res, &classification)
		if err != nil {
			return nil, err
		}
		doc.Type = param.Model.Name
		docs = append(docs, doc)
	}

	return []byte{}, nil
}

func (S *SqlDriver) QueryMultiDocumentOfProject(ctx context.Context, param shared.CommonSystemParams) ([]*shared.DefaultDocumentStructure, error) {

	var local string
	if val, ok := param.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	query, err := RootResolverQueryBuilder(param, false)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = S.ORM.NewRaw(*query).Scan(ctx, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 { // send empty response for table
		return []*shared.DefaultDocumentStructure{}, nil
	}

	classification := FieldClassification{}
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			classification.MultilineFields = append(classification.MultilineFields, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && f.Validation.IsGallery {
			classification.GalleryField = append(classification.GalleryField, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && !f.Validation.IsGallery {
			classification.PictureField = append(classification.PictureField, f.Identifier)
		} else if f.FieldType == "number" && f.InputType == "double" {
			classification.DoubleFields = append(classification.DoubleFields, f.Identifier)
		} else if f.FieldType == "list" && f.Validation != nil && (len(f.Validation.FixedListElements) == 0 || f.Validation.IsMultiChoice) {
			classification.ListFields = append(classification.ListFields, f.Identifier)
		}
	}

	var docs []*shared.DefaultDocumentStructure
	for _, res := range result {
		doc, err := CommonDocTransformation(param.Model, local, res, &classification)
		if err != nil {
			return nil, err
		}
		doc.Type = param.Model.Name
		docs = append(docs, doc)
	}

	return docs, nil
}

func (S *SqlDriver) NewInsertableRelations(ctx context.Context, param *shared.ConnectDisconnectParam) ([]string, error) {
	panic("new insertable relations not implemented")
}

func (S *SqlDriver) CheckOneToOneRelationExists(ctx context.Context, param *shared.ConnectDisconnectParam) (bool, error) {
	panic("check one to one relation not implemented")
}

func (S *SqlDriver) GetRelationIds(ctx context.Context, param *shared.ConnectDisconnectParam) ([]string, error) {
	panic("get relations ids not implemented")
}
