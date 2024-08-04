package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	_const "github.com/apito-io/databasedriver"
	"github.com/google/uuid"
	"github.com/jinzhu/inflection"
	"github.com/uptrace/bun"
	"gorm.io/datatypes"
)

func (S *SqlDriver) DeleteProject(ctx context.Context, projectId string) error {
	_, err := S.ORM.Exec(`drop schema ?;`, bun.Ident(projectId))
	if err != nil {
		return err
	}
	return nil
}

func (S *SqlDriver) DeleteMediaFile(ctx context.Context, param shared.CommonSystemParams) error {
	panic("delete media file not implemented")
}

func (S *SqlDriver) DeleteDocumentRelation(ctx context.Context, param shared.CommonSystemParams) error {
	panic("delete document relation not implemented")
}

func (S *SqlDriver) DropField(ctx context.Context, param shared.CommonSystemParams) error {

	tableName := inflection.Plural(param.Model.Name)
	_, err := S.ORM.Exec(`alter table ? drop column ?;`, bun.Ident(tableName), bun.Ident(param.FieldInfo.Identifier))
	//return nil
	if err != nil {
		return err
	}
	return nil
}

func (S *SqlDriver) RenameModel(ctx context.Context, project *protobuff.Project, modelName, newName string) error {
	panic("rename model not implemented")
}

func (S *SqlDriver) ConvertModel(ctx context.Context, project *protobuff.Project, modelName string) error {
	panic("rename model not implemented")
}

func (S *SqlDriver) RenameField(ctx context.Context, oldFieldName string, repeatedGroupIdentifier *string, param shared.CommonSystemParams) error {
	panic("rename field not implemented")
}

func (S *SqlDriver) DeleteDocumentsFromProject(ctx context.Context, param shared.CommonSystemParams) error {
	panic("delete documents from project not implemented")
}

func (S *SqlDriver) RemoveAuthAddOns(ctx context.Context, project *protobuff.Project, option map[string]interface{}) error {
	return nil
}

func (S *SqlDriver) TransferProject(ctx context.Context, userId, from, to string) error {
	return nil
}

func (S *SqlDriver) AddCollection(ctx context.Context, projectId string) (*string, error) {

	tx, err := S.ORM.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	switch S.DriverCredential.Engine {
	case _const.PostgresSQLDriver:
		// create database can not be executed inside transaction, so it's outside the transaction
		if _, err = tx.Exec("CREATE DATABASE ?", bun.Ident(projectId)); err != nil {
			return nil, err
		}

		// reinit the GORM connection
		/*S.DriverCredential.Database = projectId
		db, err := GetSQLDriver(S.DriverCredential)
		if err != nil {
			return nil, err
		}*/

		if _, err = tx.Exec(fmt.Sprintf(`
			CREATE TABLE public.meta(
				id VARCHAR(36) NOT NULL PRIMARY KEY,
				doc_id VARCHAR(36) NOT NULL,
				created_at DATE NOT NULL DEFAULT CURRENT_DATE,
				updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
				created_by VARCHAR(36) NOT NULL,
				updated_by VARCHAR(36),
				status VARCHAR(36)
			);`)); err != nil {
			return nil, err
		}

		if _, err = tx.Exec(fmt.Sprintf(`
			CREATE TABLE public.media(
				id VARCHAR(36) NOT NULL PRIMARY KEY,
				model VARCHAR(125),
				media_type VARCHAR(65),
				file_extension VARCHAR(65),
				file_name TEXT,
				size INTEGER,
				s3_key TEXT,
				url TEXT,
				created_at DATE NOT NULL DEFAULT CURRENT_DATE
			);`)); err != nil {
			return nil, err
		}

	case _const.MySQLDriver:

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if _, err = tx.Exec("CREATE DATABASE ?", bun.Ident(projectId)); err != nil {
			// return any error will rollback
			return nil, err
		}

		if _, err = tx.Exec(`
			CREATE TABLE ?.meta(
				id VARCHAR(36) NOT NULL PRIMARY KEY,
				doc_id VARCHAR(36) NOT NULL,
				created_at DATE NOT NULL DEFAULT (CURRENT_DATE),
				updated_at DATE NOT NULL DEFAULT (CURRENT_DATE),
				created_by VARCHAR(36) NOT NULL,
				updated_by VARCHAR(36),
				status VARCHAR(35)
			);`, bun.Ident(projectId)); err != nil {
			return nil, err
		}

		if _, err = tx.Exec(`
			CREATE TABLE ?.media(
				id VARCHAR(36) NOT NULL PRIMARY KEY,
				model VARCHAR(125),
				media_type VARCHAR(65),
				file_extension VARCHAR(65),
				file_name TEXT,
				size INTEGER,
				s3_key TEXT,
				url TEXT,
				created_at DATE NOT NULL DEFAULT (CURRENT_DATE)
			);`, bun.Ident(projectId)); err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		// rollback ? tx.rollback() ?
		return nil, err
	}

	return &projectId, nil
}

func (S *SqlDriver) DuplicateModel(ctx context.Context, project *protobuff.Project, modelName, newName string) (*protobuff.ProjectSchema, error) {
	//TODO implement me
	panic("duplicate model not implemented")
}

func (S *SqlDriver) AddModel(ctx context.Context, project *protobuff.Project, name string, singleRecord bool) (*protobuff.ProjectSchema, error) {

	modelType := &protobuff.ModelType{
		Name: name,
	}

	if singleRecord {
		uid := uuid.New()
		modelType.SinglePage = true
		modelType.SinglePageUuid = uid.String()
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
			return nil, errors.New("model already defined")
		}
	}

	name = inflection.Plural(name)

	//Then execute your query for creating table
	_, err := S.ORM.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS ? ( 
    		id VARCHAR(36) NOT NULL PRIMARY KEY 
        );`, bun.Ident(name))
	if err != nil {
		return nil, err
	}

	return project.Schema, nil
}

func (S *SqlDriver) AddRelationFields(ctx context.Context, from *protobuff.ConnectionType, to *protobuff.ConnectionType) error {

	toTableName := inflection.Plural(from.Model)
	fromTableName := inflection.Plural(to.Model)
	toFieldName := to.Model
	fromFieldName := from.Model

	// from connection
	switch from.Relation {
	case "has_one":
		switch to.Relation {
		case "has_one":
			//same for one to one & one to many
			_toFieldName := fmt.Sprintf(`%s_id`, toFieldName)
			_, err := S.ORM.Exec(`ALTER TABLE ? ADD ? VARCHAR(36) references ? (id) ON DELETE CASCADE;`, bun.Ident(toTableName), _toFieldName, fromTableName)
			if err != nil {
				return err
			}
			break
		case "has_many":
			//same for one to one & one to many
			_fromFieldName := fmt.Sprintf(`%s_id`, fromFieldName)
			_, err := S.ORM.Exec(`ALTER TABLE ? ADD ? VARCHAR(36) references ? (id) ON DELETE CASCADE;`, bun.Ident(fromTableName), _fromFieldName, toTableName)
			if err != nil {
				return err
			}
			break
		}
		break
	case "has_many":
		switch to.Relation {
		case "has_many":
			//same for one to one & one to many
			query := fmt.Sprintf(`CREATE TABLE %s_%s(
				%s_id VARCHAR(36) REFERENCES %s (id) ON DELETE CASCADE,
				%s_id VARCHAR(36) REFERENCES %s (id) ON DELETE CASCADE,
				PRIMARY KEY (%s_id, %s_id)
			);`, fromTableName, toTableName,
				toFieldName, fromTableName,
				fromFieldName, toTableName,
				toFieldName, fromFieldName)
			_, err := S.ORM.Exec(query)
			if err != nil {
				return err
			}
			break
		case "has_one":
			//same for one to one & one to many
			_, err := S.ORM.Exec("ALTER TABLE ? ADD ?_id VARCHAR(36) REFERENCES ? (id) ON DELETE CASCADE;",
				bun.Ident(toTableName),
				bun.Ident(toFieldName),
				bun.Ident(fromTableName),
			)
			if err != nil {
				return err
			}
			break
		}
		break
	}

	return nil
}

func (S *SqlDriver) DropConnections(ctx context.Context, projectId string, from *protobuff.ConnectionType, to *protobuff.ConnectionType) error {

	toTableName := inflection.Plural(from.Model)
	fromTableName := inflection.Plural(to.Model)
	toFieldName := to.Model
	fromFieldName := from.Model

	// from connection
	switch from.Relation {
	case "has_one":
		switch to.Relation {
		case "has_one":
			//same for one to one & one to many
			_, err := S.ORM.Exec("ALTER TABLE ? DROP CONSTRAINT ?;", bun.Ident(toTableName), toFieldName+"_id")
			if err != nil {
				return err
			}
			break
		case "has_many":
			//same for one to one & one to many
			_, err := S.ORM.Exec("ALTER TABLE ? DROP CONSTRAINT ?;", bun.Ident(fromTableName), fromFieldName+"_id")
			if err != nil {
				return err
			}
			break
		}
		break
	case "has_many":
		switch to.Relation {
		case "has_many":
			//same for one to one & one to many
			_, err := S.ORM.Exec(`DROP TABLE ?;`, fromTableName+"_"+toTableName)
			if err != nil {
				return err
			}
			break
		case "has_one":
			//same for one to one & one to many
			_, err := S.ORM.Exec("ALTER TABLE ? DROP CONSTRAINT ?;", bun.Ident(toTableName), toFieldName+"_id")
			if err != nil {
				return err
			}
			break
		}
		break
	}
	return nil
}

func (S *SqlDriver) AddFieldToModel(ctx context.Context, param shared.CommonSystemParams, isUpdate bool, repeatedGroupIdentifier *string) (*protobuff.ModelType, error) {

	if param.FieldInfo.InputType == "geo" {
		return nil, errors.New("geo Field is currently not supported in PostgresSQL. We will be integrating it soon via extension")
	}

	if repeatedGroupIdentifier == nil && !isUpdate {
		if param.FieldInfo.Serial == 0 || len(param.Model.Fields) == 0 { // new field cant be zero
			param.FieldInfo.Serial = uint32(len(param.Model.Fields) + 1)
		}
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

	if repeatedGroupIdentifier != nil {
		return param.Model, nil // dont create anything
		// todo transform this to one to many relation
	}

	var datatype string
	var validations []string
	switch param.FieldInfo.FieldType {
	case "text":
		datatype = "TEXT"
		break
	case "multiline":
		datatype = "TEXT"
		break
	case "date":
		datatype = "DATE"
		break
	case "boolean":
		datatype = "BOOLEAN"
		break
	case "media":
		if param.FieldInfo.Validation.IsGallery {
			datatype = "JSON"
		} else {
			datatype = "JSON"
		}
	case "number":
		switch param.FieldInfo.InputType {
		case "int":
			datatype = "INTEGER"
			break
		case "double":
			datatype = "NUMERIC"
			break
		}
	case "list":
		if param.FieldInfo.Validation != nil && len(param.FieldInfo.Validation.FixedListElements) > 0 && param.FieldInfo.Validation.IsMultiChoice == false {
			datatype = "TEXT"
		} else {
			datatype = "JSON"
		}
		break
	case "repeated":
		datatype = "JSON"
	}

	if param.FieldInfo.Validation != nil && param.FieldInfo.Validation.Required {
		var defaultValue interface{}
		switch param.FieldInfo.InputType {
		case "string":
			defaultValue = "''"
			break
		case "int":
			defaultValue = 0
			break
		case "bool":
			defaultValue = false
			break
		case "double":
			defaultValue = 0.0
			break
		}
		validations = append(validations, fmt.Sprintf("NOT NULL DEFAULT %v", defaultValue))
	} else if param.FieldInfo.Validation != nil && param.FieldInfo.Validation.Unique {
		validations = append(validations, "UNIQUE")
	}

	tableName := inflection.Plural(param.Model.Name)

	// local support
	if param.FieldInfo.Validation != nil && len(param.FieldInfo.Validation.Locals) > 0 {
		for _, local := range param.FieldInfo.Validation.Locals {
			var column string
			if local != "en" {
				column = fmt.Sprintf(`%s_%s`, param.FieldInfo.Identifier, local)
			} else {
				column = fmt.Sprintf(`%s`, param.FieldInfo.Identifier)
			}
			//Then execute your query for creating table
			query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s %s;", tableName, column, datatype, strings.Join(validations, " "))
			_, err := S.ORM.ExecContext(ctx, query)
			if err != nil {
				return nil, err
			}
		}
	} else {
		//Then execute your query for creating table
		query := fmt.Sprintf("ALTER TABLE %s ADD %s %s %s;", tableName, param.FieldInfo.Identifier, datatype, strings.Join(validations, " "))
		_, err := S.ORM.ExecContext(ctx, query)
		if err != nil {
			return nil, err
		}
	}

	return param.Model, nil
}

func (S *SqlDriver) AddTeamMetaInfo(ctx context.Context, docs []*protobuff.SystemUser) ([]*protobuff.SystemUser, error) {
	panic("add team meta info not implemented")
}

func (S *SqlDriver) AddATeamMemberToProject(ctx context.Context, projectId string, memberData map[string]interface{}) error {
	panic("add team member to project not implemented")
}

func (S *SqlDriver) RemoveATeamMemberFromProject(ctx context.Context, projectId string, memberId string) error {
	panic("remove a team member from project not implemented")
}

func (S *SqlDriver) CreateMediaDocument(ctx context.Context, projectId string, media *protobuff.FileDetails) (*protobuff.FileDetails, error) {

	data := map[string]interface{}{
		"id": media.Id,
	}
	if media.UploadParam != nil {
		if media.UploadParam.ModelName != "" {
			data["model"] = media.UploadParam.ModelName
		}
		if media.UploadParam.FieldName != "" {
			data["field"] = media.UploadParam.FieldName
		}
	}
	if media.ContentType != "" {
		data["media_type"] = media.ContentType
	}
	if media.FileExtension != "" {
		data["file_extension"] = media.FileExtension
	}
	if media.FileName != "" {
		data["file_name"] = media.FileName
	}
	if media.Size != 0 {
		data["size"] = media.Size
	}
	if media.S3Key != "" {
		data["s3_key"] = media.S3Key
	}
	if media.Url != "" {
		data["url"] = media.Url
	}

	_, err := S.ORM.NewInsert().Table("media").Model(data).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (S *SqlDriver) AddDocumentToProject(ctx context.Context, projectId string, modelName string, doc *shared.DefaultDocumentStructure) (interface{}, error) {

	tx, err := S.ORM.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	tableName := inflection.Plural(modelName)

	data := map[string]interface{}{
		"id": doc.Id,
	}
	for k, v := range doc.Data {
		if val, ok := v.(map[string]interface{}); ok {
			if html, ok := val["html"]; ok {
				data[k] = html
			}
		} else {
			data[k] = v
		}
	}
	_, err = tx.NewInsert().Table(tableName).Model(data).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// now insert a meta data
	metaData := map[string]interface{}{
		"id":         uuid.New().String(),
		"created_by": doc.Meta.CreatedBy.Id,
		"updated_by": doc.Meta.LastModifiedBy.Id,
		"status":     doc.Meta.Status,
		"doc_id":     doc.Id,
	}
	_, err = tx.NewInsert().Table("meta").Model(metaData).Exec(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (S *SqlDriver) UpdateDocumentOfProject(ctx context.Context, param shared.CommonSystemParams, doc *shared.DefaultDocumentStructure, replace bool) error {

	var multilineFields []string
	var pictureField []string
	var galleryField []string
	var listFields []string
	var repeatedFields []string
	for _, f := range param.Model.Fields {
		if f.FieldType == "multiline" {
			multilineFields = append(multilineFields, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && f.Validation.IsGallery {
			galleryField = append(galleryField, f.Identifier)
		} else if f.FieldType == "media" && f.Validation != nil && !f.Validation.IsGallery {
			pictureField = append(pictureField, f.Identifier)
		} else if f.FieldType == "list" && f.Validation != nil && (len(f.Validation.FixedListElements) == 0 || f.Validation.IsMultiChoice) {
			listFields = append(listFields, f.Identifier)
		} else if f.FieldType == "repeated" {
			repeatedFields = append(repeatedFields, f.Identifier)
		}
	}

	tableName := inflection.Plural(doc.Type)

	data := map[string]interface{}{}
	for k, v := range doc.Data {
		// if its a map then it must be a media field
		kind := reflect.ValueOf(v).Kind()
		switch kind {
		case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
			data[k] = v
			break
		case reflect.Map:
			val := v.(map[string]interface{})
			if utility.ArrayContains(multilineFields, k) {
				if html, ok := val["html"]; ok {
					data[k] = html
				}
			} else if utility.ArrayContains(pictureField, k) {
				b, _ := json.Marshal(v)
				data[k] = datatypes.JSON(b)
			}
			break
		case reflect.Ptr:
			fmt.Println(v)
			break
		case reflect.Slice:
			if utility.ArrayContains(galleryField, k) || utility.ArrayContains(listFields, k) || utility.ArrayContains(repeatedFields, k) {
				b, err := json.Marshal(v)
				if err != nil {
					return err
				}
				data[k] = datatypes.JSON(b)
			}
			break
		}

	}
	_, err := S.ORM.NewUpdate().Table(tableName).Where("id = ?", doc.Id).Model(data).Exec(ctx)
	if err != nil {
		return err
	}

	// now insert a meta data
	metaData := map[string]interface{}{
		"updated_at": utility.GetCurrentTime(),
		"updated_by": param.UserId,
	}

	_, err = S.ORM.NewUpdate().Table("meta").Where("doc_id = ?", doc.Id).Model(metaData).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (S *SqlDriver) DeleteDocumentFromProject(ctx context.Context, param shared.CommonSystemParams) error {

	tableName := inflection.Plural(param.Model.Name)

	_, err := S.ORM.NewDelete().Table(tableName).Where("id = ?", param.DocumentId).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (S *SqlDriver) CreateRelation(ctx context.Context, projectId string, relation *shared.EdgeRelation) error {
	panic("create relation not implemented")
}

func (S *SqlDriver) DeleteRelation(ctx context.Context, param *shared.ConnectDisconnectParam, id string) error {
	panic("delete relation not implemented")
}
