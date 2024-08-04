package sql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"

	"github.com/graph-gophers/dataloader"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jinzhu/inflection"
	"github.com/tailor-inc/graphql"
)

type QueryBuilderParam struct {
	CollectionName string
	RelationName   string
	Args           map[string]interface{}
	ParentModel    string
	ModelType      *protobuff.ModelType
}

var FilterSuffix = map[string]string{
	"eq":     "==",
	"ne":     "!=",
	"lt":     "<",
	"lte":    "<=",
	"gt":     ">",
	"gte":    ">=",
	"in":     "IN",
	"not_in": "NOT IN",
}

func SelectBuilder(mv string, local string, modelType *protobuff.ModelType, returnCount bool) []string {
	var returnType []string
	metaQuery := fmt.Sprintf(`%s.created_at AS sys_created_at, %s.updated_at AS sys_updated_at, %s.created_by AS sys_created_by, %s.updated_by AS sys_updated_by, %s.status as sys_status`, mv, mv, mv, mv, mv)
	if local != "" {
		var dataJson []string
		for _, f := range modelType.Fields {
			if f.Validation != nil && local != "en" && utility.ArrayContains(f.Validation.Locals, local) {
				dataJson = append(dataJson, fmt.Sprintf(`x.%s_%s AS %s`, f.Identifier, local, f.Identifier))
			} else {
				dataJson = append(dataJson, fmt.Sprintf(`x.%s`, f.Identifier))
			}
		}
		returnType = append(returnType, []string{"x.id", strings.Join(dataJson, ", "), metaQuery}...)
	} else if returnCount {
		returnType = append(returnType, "count(x.id)")
	} else {
		returnType = append(returnType, []string{"x.*", metaQuery}...)
	}
	return returnType
}

func LimitBuilder(param *graphql.ResolveParams) (int, int) {
	arg := param.Args

	limit := 10
	if val, ok := arg["limit"]; ok {
		limit = val.(int)
	}

	start := 0
	if val, ok := arg["start"]; ok {
		start = val.(int)
	}

	page := 1
	if val, ok := arg["page"]; ok {
		page = val.(int)
	}

	if page > 1 {
		offset := limit * (page - 1)
		return limit, offset
	}
	return limit, start
}

func getFieldType(val interface{}) reflect.Kind {
	return reflect.TypeOf(val).Kind()
}

func getUserFieldType(fieldInfo *protobuff.FieldInfo) reflect.Kind {
	switch fieldInfo.InputType {
	case "string":
		switch fieldInfo.FieldType {
		case "list":
			if !fieldInfo.Validation.IsMultiChoice && len(fieldInfo.Validation.FixedListElements) > 0 { // and multi-choice & dynamic list
				return reflect.String
			} else {
				return reflect.Slice
			}
		case "media":
		case "multiline":
		default:
			return reflect.String
		}
		break
	case "int":
		return reflect.Int
		break
	case "double":
		return reflect.Float64
		break
	case "bool":
		return reflect.Bool
		break
	}
	return reflect.Interface
}

func ConditionBuilder(variable string, args map[string]interface{}, modelType *protobuff.ModelType) (map[string][]string, error) {

	userDefinedFieldNames := make(map[string]reflect.Kind)
	for _, field := range modelType.Fields {
		userDefinedFieldNames[field.Identifier] = getUserFieldType(field)
	}

	conditions := make(map[string][]string)

	if w, ok := args["where"]; w != nil && ok {
		where := w.(map[string]interface{})
		for field, filterObj := range where {

			switch field { // default is x.data
			case "role":
				variable = "x"
			case "tenant":
				variable = "x"
			}

			// if AND / OR found
			switch field {
			case "AND":
				conditions["AND"], _ = FilterBuilder(variable, filterObj.(map[string]interface{}), modelType)
			case "OR":
				conditions["OR"], _ = FilterBuilder(variable, filterObj.(map[string]interface{}), modelType)
			default:
				conditions["AND"], _ = FilterBuilder(variable, where, modelType)
			}
		}
	}

	return conditions, nil
}

func FilterBuilder(variable string, where map[string]interface{}, modelType *protobuff.ModelType) ([]string, error) {

	userDefinedFieldNames := make(map[string]reflect.Kind)
	for _, field := range modelType.Fields {
		userDefinedFieldNames[field.Identifier] = getUserFieldType(field)
	}

	var conditions []string

	for field, filterObj := range where {

		fieldName := fmt.Sprintf("%s.%s", variable, field)

		var actualValue interface{}
		for suffix, value := range filterObj.(map[string]interface{}) {

			actualValue = value

			switch suffix {
			case "contains":
				conditions = append(conditions, fmt.Sprintf(`%s LIKE '%%%s%%'`, fieldName, value.(string)))
				break
			case "eq", "ne", "lt", "lte", "gt", "gtr", "in", "not_in":
				switch value.(type) {
				case int, float64, bool:
					conditions = append(conditions, fmt.Sprintf(`%s %s %v`, fieldName, FilterSuffix[suffix], value.(string)))
					break
				case string:
					conditions = append(conditions, fmt.Sprintf(`%s %s '%v'`, fieldName, FilterSuffix[suffix], value.(string)))
					break
				case []interface{}:
					var vals []string
					for _, v := range value.([]interface{}) {
						switch v.(type) {
						case int, float64:
							vals = append(vals, fmt.Sprintf("%v", v))
						case string:
							vals = append(vals, fmt.Sprintf("'%v'", v))
						}
					}
					final := fmt.Sprintf("[%s]", strings.Join(vals, ","))
					conditions = append(conditions, fmt.Sprintf(`COUNT(%s[* FILTER CONTAINS(%s, CURRENT)])`, fieldName, final))
				}
				break
			}
		}

		//validate the field & type
		if kind, ok := userDefinedFieldNames[field]; ok {
			k := getFieldType(actualValue)
			if kind != k {
				return nil, errors.New(fmt.Sprintf("Invalid Value for %s in Query. Type mismatched", field))
			}
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid Field Name %s in Query", field))
		}

	}

	return conditions, nil
}

func CommonDocTransformation(model *protobuff.ModelType, local string, result map[string]interface{}, classification *FieldClassification) (*shared.DefaultDocumentStructure, error) {

	doc := shared.DefaultDocumentStructure{
		Type: model.Name,
		Meta: &protobuff.MetaField{},
	}

	if val, ok := result["id"].(string); ok {
		doc.Id = string(val)
		doc.Key = doc.Id
	} else {
		return nil, errors.New("id is required for any document to fetch")
	}

	data := map[string]interface{}{}

	for k, v := range result {
		switch k {
		case "doc_id":
			continue
		case "sys_key":
			continue
		case "id":
			continue
			break
		case "sys_status":
			doc.Meta.Status = v.(string)
			break
		case "sys_created_at":
			t := time.Unix(v.(time.Time).Unix(), 0)
			doc.Meta.CreatedAt = t.Format(time.RFC3339)
			break
		case "sys_updated_at":
			t := time.Unix(v.(time.Time).Unix(), 0)
			doc.Meta.UpdatedAt = t.Format(time.RFC3339)
			break
		case "sys_updated_by":
			id := v.(string)
			doc.Meta.LastModifiedBy = &protobuff.SystemUser{
				Id: string(id),
			}
			break
		case "sys_created_by":
			id := v.(string)
			doc.Meta.CreatedBy = &protobuff.SystemUser{
				Id: string(id),
			}
			break
		default:
			if utility.ArrayContains(classification.MultilineFields, k) {
				converter := md.NewConverter("", true, nil)
				var html string
				if val, ok := v.(string); ok {
					html = val
				}
				markdown, err := converter.ConvertString(html)
				if err != nil {
					fmt.Println(err.Error())
				}
				data[k] = map[string]interface{}{
					"html":     html,
					"markdown": markdown,
					"text":     strip.StripTags(html),
				}
			} else if utility.ArrayContains(classification.DoubleFields, k) {
				if val, ok := v.([]byte); ok {
					f, _ := strconv.ParseFloat(string(val), 64)
					data[k] = f
				}
			} else if utility.ArrayContains(classification.PictureField, k) {
				if val, ok := v.([]byte); ok {
					var pic map[string]interface{}
					err := json.Unmarshal(val, &pic)
					if err != nil {
						return nil, err
					}
					data[k] = pic
				}
			} else if utility.ArrayContains(classification.GalleryField, k) {
				if val, ok := v.([]byte); ok {
					var gallery []map[string]interface{}
					err := json.Unmarshal(val, &gallery)
					if err != nil {
						return nil, err
					}
					data[k] = gallery
				}
			} else if utility.ArrayContains(classification.ListFields, k) {
				if val, ok := v.([]byte); ok {
					var lists []interface{}
					err := json.Unmarshal(val, &lists)
					if err != nil {
						return nil, err
					}
					data[k] = lists
				}
			} else if subfields, ok := classification.RepeatedFields[k]; ok && len(classification.RepeatedFields) > 0 {

				var repeated []map[string]interface{}
				if val, ok := v.([]byte); ok {
					err := json.Unmarshal(val, &repeated)
					if err != nil {
						return nil, err
					}
				}
				if local == "en" {
					data[k] = repeated
				} else {
					for _, subItem := range repeated {
						for _, f := range subfields {
							if f.Validation != nil && utility.ArrayContains(f.Validation.Locals, local) {
								if localContentFound, ok := subItem[fmt.Sprintf(`%s_%s`, f.Identifier, local)]; ok {
									subItem[f.Identifier] = localContentFound
								}
								break
							}
						}
					}
					data[k] = repeated
				}
			} else {
				data[k] = v
			}
		}
	}
	doc.Data = data
	return &doc, nil
}

func MediaDocTransformation(docType string, result map[string]interface{}) (*protobuff.FileDetails, error) {

	doc := protobuff.FileDetails{
		Type: docType,
	}

	if val, ok := result["id"].([]byte); ok {
		doc.Id = string(val)
		doc.XKey = doc.Id
	} else {
		return nil, nil
	}

	if val, ok := result["created_at"].(time.Time); ok {
		t := time.Unix((val).Unix(), 0)
		doc.CreatedAt = t.Format(time.RFC3339)
	}

	if val, ok := result["model"].(string); ok {
		if doc.UploadParam == nil {
			doc.UploadParam = &protobuff.UploadParams{}
		}
		doc.UploadParam.ModelName = val
	}

	if val, ok := result["s3_key"].(string); ok {
		doc.S3Key = val
	}

	if val, ok := result["media_type"].(string); ok {
		doc.ContentType = val
	}

	if val, ok := result["file_extension"].(string); ok {
		doc.FileExtension = val
	}

	if val, ok := result["file_name"].(string); ok {
		doc.FileName = val
	}

	if val, ok := result["size"].(int32); ok {
		doc.Size = int64(val)
	}

	if val, ok := result["url"].(string); ok {
		doc.Url = val
	}

	return &doc, nil
}

func RootConnectionResolverQueryBuilder(param *shared.CommonSystemParams) (string, error) {

	projectId := param.ProjectId
	_args := param.ResolveParams.Args

	filters, err := ConditionBuilder("x.data", _args, param.Model)
	if err != nil {
		return "", err
	}

	var mergedFilter []string
	for condition, _ := range filters {
		mergedFilter = append(mergedFilter, strings.Join(filters[condition], condition))
	}

	model := param.Model.Name

	var queries []string
	queries = append(queries, fmt.Sprintf("FOR x in `p_%s`", projectId))
	if len(filters) > 0 {
		queries = append(queries, fmt.Sprintf(`Filter x.type == '%s' AND %s`, model, strings.Join(mergedFilter, " AND "))) // #todo need fix
	} else {
		queries = append(queries, fmt.Sprintf(`Filter x.type == '%s'`, model))
	}
	queries = append(queries, fmt.Sprintf(`COLLECT WITH COUNT INTO total`))
	queries = append(queries, fmt.Sprintf(`return total`))

	return strings.Join(queries, " "), nil
}

func RootResolverQueryBuilder(param shared.CommonSystemParams, returnCount bool) (*string, error) {

	/*var modelName string
	if param.Model == nil {
		modelName = inflection.Singular(param.ResolveParams.Info.FieldName)
	} else {
		modelName = inflection.Singular(param.Model.Name)
	}*/

	var leftJoins []string

	var local string
	if val, ok := param.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	returnType := SelectBuilder("y", local, param.Model, returnCount)

	var connection map[string]interface{}
	if val, ok := param.ResolveParams.Args["connection"].(map[string]interface{}); len(val) > 0 && ok {
		connection = val

		var connectionType string
		if val, ok := connection["connection_type"].(string); ok && val != "" {
			connectionType = val
		} else {
			return nil, errors.New("connection type is required if passing connection object")
		}

		var fromModel string
		var toModel string
		switch connectionType {
		case "forward":
			fromModel = connection["to_model"].(string)
			toModel = connection["model"].(string)
		case "backward":
			fromModel = connection["model"].(string)
			toModel = connection["to_model"].(string)
		default:
			return nil, errors.New("invalid connection type")
		}
		pivotTable := fmt.Sprintf(`%s_%s`, inflection.Plural(fromModel), inflection.Plural(toModel))
		leftJoins = append(leftJoins, fmt.Sprintf(`left join %s as z on z.%s_id = x.id`, pivotTable, toModel))

		// z.* has no use in return count
		if !returnCount {
			returnType = append(returnType, "z.*")
		} else {
			toModel = connection["model"].(string)
			returnType = []string{fmt.Sprintf("z.%s_id", toModel)}
		}
	}

	limit, offset := LimitBuilder(param.ResolveParams)

	tableName := inflection.Plural(param.Model.Name)

	if !returnCount {
		leftJoins = append(leftJoins, fmt.Sprintf(`left join meta as y on y.doc_id = x.id`))
	}

	var queries []string
	queries = append(queries, fmt.Sprintf(`SELECT %s FROM %s as x %s`,
		strings.Join(returnType, ", "),
		tableName,
		strings.Join(leftJoins, "\n"),
	))

	filters, err := ConditionBuilder("x", param.ResolveParams.Args, param.Model)
	if err != nil {
		return nil, err
	}

	// filter based on roles
	// #todo fix this
	/*if permission, ok := param.Role.ApiPermissions[modelName]; ok && param.Role.ApiPermissions != nil {
		switch permission.Read {
		case "own":
			filters["AND"] = []string{fmt.Sprintf(`y.created_by == '%s'`, param.UserId)}
			break
		case "tenant":
			filters["AND"] = []string{fmt.Sprintf(`y.tenant_id == '%s'`, param.TenantId)}
			break
		}
	}*/

	var mergedFilter []string
	for condition, _ := range filters {
		mergedFilter = append(mergedFilter, fmt.Sprintf(`(%s)`, strings.Join(filters[condition], fmt.Sprintf(` %s `, condition))))
	}

	var intersect bool
	if val, ok := param.ResolveParams.Args["intersect"].(bool); ok {
		intersect = val
	}

	if connection != nil && len(connection) > 0 {
		fromModel := connection["to_model"].(string)
		toModel := connection["model"].(string)
		id := connection["_id"].(string)

		if intersect { // get the other result
			subQuery := fmt.Sprintf(`SELECT z.%s_id FROM %s as x %s WHERE z.%s_id = '%s'`,
				toModel,
				inflection.Plural(toModel),
				strings.Join(leftJoins, "\n"),
				fromModel,
				id,
			)
			mergedFilter = append(mergedFilter, fmt.Sprintf(`x.id not in (%s)`, subQuery))
		} else { // else get the exact match
			mergedFilter = append(mergedFilter, fmt.Sprintf(`z.%s_id = '%s'`, fromModel, id))
		}
	}

	if len(filters) > 0 || connection != nil {
		queries = append(queries, fmt.Sprintf(`WHERE %s`, strings.Join(mergedFilter, " AND ")))
	}

	// default sort
	if !returnCount { // limit & Order if not counting
		queries = append(queries, fmt.Sprintf(`ORDER BY y.created_at DESC`))
		queries = append(queries, fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset))
	}

	query := strings.Join(queries, " ")

	return &query, nil
}

func BuildCombinedRelationQuery(relationType string, parentModel string, arg *shared.CommonSystemParams) (*string, *string, error) {

	var local string
	if val, ok := arg.ResolveParams.Args["local"].(string); ok {
		local = val
	}

	filters, err := ConditionBuilder("x", arg.ResolveParams.Args, arg.Model)
	if err != nil {
		return nil, nil, err
	}

	var mergedFilter []string
	for condition, _ := range filters {
		mergedFilter = append(mergedFilter, strings.Join(filters[condition], condition))
	}

	var relationshipDirection string
	for _, m := range arg.Model.Connections {
		if m.Model == parentModel && m.Type == "backward" {
			relationshipDirection = "to"
		} else if m.Model == parentModel && m.Type == "forward" {
			relationshipDirection = "from"
		}
	}

	if relationshipDirection == "" {
		return nil, nil, errors.New("Could not decide form/to relations")
	}

	relationInput := map[string]interface{}{}
	if len(arg.ResolveParams.Args) > 0 {
		relationInput = arg.ResolveParams.Args
	} else {
		relationInput = map[string]interface{}{
			"from_model":    parentModel,
			"to_model":      arg.Model.Name,
			"relation_type": relationType,
		}
	}

	keys := arg.DocumentIDs

	selectThing := SelectBuilder("z", local, arg.Model, arg.OnlyReturnCount)

	relationTo := relationInput["to_model"].(string)
	tableName := inflection.Plural(relationTo)

	var query string
	var whereCondition string
	var pivotTable string
	switch relationInput["relation_type"] {
	case "has_many":
		var manyToManyRelation bool
		for _, c := range arg.Model.Connections {
			if c.Model == parentModel {
				if c.Relation == "has_many" {
					manyToManyRelation = true
					break
				}
			}
		}

		var keyField string
		if manyToManyRelation {

			if len(filters) > 0 {
				whereCondition = fmt.Sprintf(`y.%s_id IN ('%s') AND %s`, parentModel, strings.Join(keys, "','"), strings.Join(mergedFilter, " AND "))
			} else {
				whereCondition = fmt.Sprintf(`y.%s_id IN ('%s')`, parentModel, strings.Join(keys, "','"))
			}

			keyField = fmt.Sprintf(`y.%s_id`, parentModel)

			switch relationshipDirection {
			case "to":
				pivotTable = fmt.Sprintf(`%s_%s`, inflection.Plural(relationInput["from_model"].(string)), inflection.Plural(relationTo))
				break
			case "from":
				pivotTable = fmt.Sprintf(`%s_%s`, inflection.Plural(relationTo), inflection.Plural(relationInput["from_model"].(string)))
				break
			}
			query = fmt.Sprintf(`SELECT %s as key, %s FROM %s AS y 
				LEFT JOIN %s AS x ON x.id = y.%s_id 
				LEFT JOIN meta AS z ON z.doc_id = x.id
				WHERE %s`, keyField, selectThing, pivotTable, tableName, relationTo, whereCondition)
		} else {

			if len(filters) > 0 {
				whereCondition = fmt.Sprintf(`y.id IN ('%s') AND %s`, strings.Join(keys, "','"), strings.Join(mergedFilter, " AND "))
			} else {
				whereCondition = fmt.Sprintf(`y.id IN ('%s')`, strings.Join(keys, "','"))
			}

			keyField = fmt.Sprintf(`y.id`)

			switch relationshipDirection {
			case "to":
				pivotTable = fmt.Sprintf(`%s`, inflection.Plural(relationInput["from_model"].(string)))
				break
			case "from":
				pivotTable = fmt.Sprintf(`%s`, inflection.Plural(relationTo))
				break
			}
			query = fmt.Sprintf(`SELECT %s AS sys_key, %s FROM %s AS y 
				LEFT JOIN %s AS x ON x.%s_id = y.id 
				LEFT JOIN meta AS z ON z.doc_id = x.id
				WHERE %s`, keyField, selectThing, pivotTable, tableName, relationInput["from_model"].(string), whereCondition)
		}
		break
	case "has_one":
		if len(filters) > 0 {
			whereCondition = fmt.Sprintf(`y.id IN ('%s') AND %s`, strings.Join(keys, "','"), strings.Join(mergedFilter, " AND "))
		} else {
			whereCondition = fmt.Sprintf(`y.id IN ('%s') `, strings.Join(keys, "','"))
		}

		keyField := fmt.Sprintf(`y.id`)

		switch relationshipDirection {
		case "to":
			pivotTable = fmt.Sprintf(`%s`, inflection.Plural(relationInput["from_model"].(string)))
			query = fmt.Sprintf(`SELECT %s AS sys_key, %s FROM %s AS y 
				LEFT JOIN %s AS x ON x.%s_id = y.id 
				LEFT JOIN meta AS z ON z.doc_id = x.id
				WHERE %s LIMIT 1`, keyField, selectThing, pivotTable, tableName, relationInput["from_model"].(string), whereCondition)
			break
		case "from":
			pivotTable = fmt.Sprintf(`%s`, inflection.Plural(relationInput["from_model"].(string)))
			query = fmt.Sprintf(`SELECT %s AS sys_key, %s FROM %s AS y 
				LEFT JOIN %s AS x ON x.id = y.%s_id 
				LEFT JOIN meta AS z ON z.doc_id = x.id
				WHERE %s LIMIT 1`, keyField, selectThing, pivotTable, tableName, relationTo, whereCondition)
			break
		}
		break
	}

	rt := relationInput["relation_type"].(string)

	return &query, &rt, nil
}

func RootResolverMediaQueryBuilder(param *graphql.ResolveParams) (string, error) {

	limit, offset := LimitBuilder(param)
	var queries []string
	queries = append(queries, fmt.Sprintf(`SELECT * FROM media AS x`))

	if val, ok := param.Args["model"].(string); ok {
		if val != "" {
			queries = append(queries, fmt.Sprintf(`WHERE x.model = '%s'`, val))
		}
	} else if val, ok := param.Args["search"]; ok {
		queries = append(queries, fmt.Sprintf(`WHERE x.file_name LIKE '%%%s%%'`, val))
	}

	// default sort
	queries = append(queries, fmt.Sprintf(`ORDER BY x.created_at DESC`))
	queries = append(queries, fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset))

	return strings.Join(queries, " "), nil
}

func BuildCombinedMetaQuery(keys dataloader.Keys, param *QueryBuilderParam) ([]byte, error) {

	queries := make(map[string]string)
	for _, key := range keys {
		meta := key.(*models.ResolverKey).GetMeta()
		metaUserIDs := []string{meta.CreatedBy.Id, meta.LastModifiedBy.Id}
		queries[key.String()] = fmt.Sprintf(`(FOR u IN users FILTER u._key in ['%s'] return u)`, strings.Join(metaUserIDs, `','`))
	}

	query, err := json.Marshal(queries)
	if err != nil {
		return nil, err
	}

	first := bytes.ReplaceAll(query, []byte(`:"(`), []byte(`:(`))
	query = bytes.ReplaceAll(first, []byte(`u)"`), []byte(`u)`))

	return bytes.Join([][]byte{[]byte(`return`), query}, []byte(" ")), nil
}
