package helper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/apito-io/databasedriver/utility"
	"github.com/vektah/gqlparser/v2/ast"
)

func AssignEmptyContent(fieldDetails *protobuff.FieldInfo) interface{} {
	switch fieldDetails.InputType {
	case "string":
		return ""
	default:
		return ""
	}
}

func SelectionToFieldBuilder(selections []ast.Selection) []*protobuff.FieldInfo {
	var fields []*protobuff.FieldInfo
	for _, s := range selections {
		f := &protobuff.FieldInfo{
			Identifier:      "Field",
			InputType:       "string",
			FieldType:       "text",
			SystemGenerated: true,
		}
		if _s := s.(*ast.Field).SelectionSet; _s != nil {
			f.SubFieldInfo = SelectionToFieldBuilder(_s)
		}
		fields = append(fields, f)
	}
	return fields
}

func FieldToSelectionBuilder(fields []*protobuff.FieldInfo) []ast.Selection {
	var sections []ast.Selection
	for _, f := range fields {
		s := &ast.Field{
			/*Kind: "Field",
			Name: &ast.Name{
				Kind:  "Name",
				Value: f.Identifier,
			},*/
			Name:             f.Identifier,
			ObjectDefinition: &ast.Definition{Kind: "Field"},
		}
		// inject for the exception
		switch f.FieldType {
		case "media":
			f.SubFieldInfo = []*protobuff.FieldInfo{
				{Identifier: "url", FieldType: "text", InputType: "string"},
				{Identifier: "id", FieldType: "text", InputType: "string"},
				{Identifier: "file_name", FieldType: "text", InputType: "string"},
			}
			break
		case "multiline":
			f.SubFieldInfo = []*protobuff.FieldInfo{
				{Identifier: "html", FieldType: "text", InputType: "string"},
				{Identifier: "markdown", FieldType: "text", InputType: "string"},
				{Identifier: "text", FieldType: "text", InputType: "string"},
			}
			break
		case "geo":
			f.SubFieldInfo = []*protobuff.FieldInfo{
				{Identifier: "coordinates", FieldType: "text", InputType: "double"},
				{Identifier: "lat", FieldType: "text", InputType: "string"},
				{Identifier: "lon", FieldType: "text", InputType: "string"},
				{Identifier: "type", FieldType: "text", InputType: "string"},
			}
			break
		}
		if f.SubFieldInfo != nil {
			//Kind:       "SelectionSet",
			s.SelectionSet = ast.SelectionSet{
				&ast.Field{
					Alias:            "",
					Name:             "",
					Arguments:        nil,
					Directives:       nil,
					SelectionSet:     FieldToSelectionBuilder(f.SubFieldInfo),
					Position:         nil,
					Comment:          nil,
					Definition:       nil,
					ObjectDefinition: &ast.Definition{Kind: "SelectionSet"},
				},
			}
		}
		sections = append(sections, s)
	}
	return sections
}

func ReturnBuilder(_var string, local string, _fields *shared.FieldDetails, set ast.SelectionSet) map[string]*shared.FieldDetails {

	fieldInfos := make(map[string]*shared.FieldDetails)
	for _, f := range _fields.SubFields {
		fd := &shared.FieldDetails{
			FieldType:  f.FieldType,
			SubFields:  f.SubFieldInfo,
			Validation: f.Validation,
		}
		if f.Validation != nil && utility.ArrayContains(f.Validation.Locals, local) {
			fd.Local = local
		}

		// type and id should be included by default in any return query
		if _fields.FieldType == "" && utility.ArrayContains([]string{"id", "type"}, f.Identifier) {
			fd.Value = fmt.Sprintf(`%s.%s`, _var, f.Identifier)
		}

		// include the html query by default
		if _fields.FieldType == "multiline" && f.Identifier == "html" {
			fd.Value = fmt.Sprintf(`%s.html`, _var)
		}

		fieldInfos[f.Identifier] = fd
	}

	for _, ss := range set {
		x := ss.(*ast.Field)

		if x.Name == "__typename" { // skip all the __typename included by client
			continue
		}

		if x.Name == "data" && x.SelectionSet == nil { // inject the selections if empty, used for `listSingleModelData` query
			var dataFields []*protobuff.FieldInfo
			for _, f := range _fields.SubFields {
				if f.Identifier == "data" {
					dataFields = f.SubFieldInfo
					break
				}
			}
			//Kind:       "SelectionSet",
			x.SelectionSet = ast.SelectionSet{
				&ast.Field{
					Alias:            "",
					Name:             "",
					Arguments:        nil,
					Directives:       nil,
					SelectionSet:     FieldToSelectionBuilder(dataFields),
					Position:         nil,
					Comment:          nil,
					Definition:       nil,
					ObjectDefinition: &ast.Definition{Kind: "SelectionSet"},
				},
			}
		}

		var fieldDetails *shared.FieldDetails
		if val, ok := fieldInfos[x.Name]; ok && val != nil {
			fieldDetails = val

			// inject system subfields for special cases
			switch val.FieldType {
			case "media":
				val.SubFields = []*protobuff.FieldInfo{
					{Identifier: "url", FieldType: "text", InputType: "string"},
					{Identifier: "id", FieldType: "text", InputType: "string"},
					{Identifier: "file_name", FieldType: "text", InputType: "string"},
				}
				break
			case "multiline":
				val.SubFields = []*protobuff.FieldInfo{
					{Identifier: "html", FieldType: "text", InputType: "string"},
					{Identifier: "markdown", FieldType: "text", InputType: "string"},
					{Identifier: "text", FieldType: "text", InputType: "string"},
				}
				break
			case "geo":
				val.SubFields = []*protobuff.FieldInfo{
					{Identifier: "coordinates", FieldType: "text", InputType: "double"},
					{Identifier: "lat", FieldType: "text", InputType: "string"},
					{Identifier: "lon", FieldType: "text", InputType: "string"},
					{Identifier: "type", FieldType: "text", InputType: "string"},
				}
				break
			}
		} else {
			continue
			fmt.Println(x)
		}

		if x.SelectionSet != nil {
			var name string
			if fieldDetails.Local != "" && fieldDetails.Local != "en" {
				name = fmt.Sprintf(`%s.%s_%s`, _var, x.Name, fieldDetails.Local)
			} else {
				name = fmt.Sprintf(`%s.%s`, _var, x.Name)
			}
			for _, f := range _fields.SubFields {
				if f.Identifier == x.Name && f.SubFieldInfo != nil {
					_fields.SubFields = f.SubFieldInfo
					break
				}
			}
			generated := ReturnBuilder(name, local, fieldDetails, x.SelectionSet)
			fieldInfos[x.Name].Value = generated
		} else {
			var name string
			if fieldDetails.Local != "" && fieldDetails.Local != "en" {
				name = fmt.Sprintf(`%s.%s_%s`, _var, x.Name, fieldDetails.Local)
			} else {
				name = fmt.Sprintf(`%s.%s`, _var, x.Name)
			}
			fieldInfos[x.Name].Value = name
		}
	}
	return fieldInfos
}

func MapApitoFieldType2(fieldInfo *protobuff.FieldInfo) *shared.FieldDetails {
	switch fieldInfo.InputType {
	case "string":
		switch fieldInfo.FieldType {
		case "list":
			if !fieldInfo.Validation.IsMultiChoice && len(fieldInfo.Validation.FixedListElements) > 0 { // dropdown
				return &shared.FieldDetails{
					Kind:      reflect.String,
					SubFields: nil,
				}
			} else {
				return &shared.FieldDetails{
					Kind:      reflect.Slice,
					SubFields: nil,
				}
			}
		case "multiline", "media":
			return &shared.FieldDetails{
				Kind:      reflect.String,
				SubFields: nil,
			}
		default:
			return &shared.FieldDetails{
				Kind:      reflect.String,
				SubFields: nil,
			}
		}
	case "int":
		return &shared.FieldDetails{
			Kind:      reflect.Int,
			SubFields: nil,
		}
	case "double":
		return &shared.FieldDetails{
			Kind:      reflect.Float64,
			SubFields: nil,
		}
	case "bool":
		return &shared.FieldDetails{
			Kind:      reflect.Bool,
			SubFields: nil,
		}
	case "geo":
		return &shared.FieldDetails{
			Kind:      reflect.Map,
			SubFields: nil,
		}
	}
	return &shared.FieldDetails{
		Kind:      reflect.Interface,
		SubFields: fieldInfo.SubFieldInfo,
	}
}

/*
	func returnAQLObjectBuilderBk(_var string, _pv string, nestedMedia bool, _map map[string]*FieldDetails) ([]string, string, error) {
		var vals []string
		for k, v := range _map {
			if v.Value != nil {
				switch v.FieldType {
				case "repeated":
					if !contains([]string{"data", "meta", "created_by", "last_modified_by"}, k) { // skip for data object
						_nestedVar := utility.RandomVariableGenerator(4)
						_returns, _pvr, err := returnAQLObjectBuilderBk(_nestedVar, _pv, false, v.Value.(map[string]*FieldDetails))
						if err != nil {
							return nil, "", err
						}
						vals = append(vals, fmt.Sprintf(`"%s" : ( FOR %s in NOT_NULL(%s) ? %s : [] RETURN { %s } )`, k, _nestedVar, _pvr, _pvr, strings.Join(_returns, ", ")))
					} else {
						_returns, _, err := returnAQLObjectBuilderBk(_var, _pv, false, v.Value.(map[string]*FieldDetails))
						if err != nil {
							return nil, "", err
						}
						vals = append(vals, fmt.Sprintf(`"%s" : { %s }`, k, strings.Join(_returns, ", ")))
					}
					break
				case "media", "multiline", "geo":
					if v.Validation != nil && v.Validation.IsGallery { // multiple media is an array
						_nestedVar := utility.RandomVariableGenerator(4)
						_returns, _pvr, err := returnAQLObjectBuilderBk(_nestedVar, _pv, false, v.Value.(map[string]*FieldDetails))
						if err != nil {
							return nil, "", err
						}
						vals = append(vals, fmt.Sprintf(`"%s" : ( FOR %s in NOT_NULL(%s) ? %s : [] RETURN { %s } )`, k, _nestedVar, _pvr, _pvr, strings.Join(_returns, ", ")))
					} else {
						_returns, _pvr, err := returnAQLObjectBuilderBk(_var, _pv, true, v.Value.(map[string]*FieldDetails))
						if err != nil {
							return nil, "", err
						}
						_pv = _pvr
						vals = append(vals, fmt.Sprintf(`"%s" : { %s }`, k, strings.Join(_returns, ", ")))
					}
					break
				default:
					if _var == "" {
						vals = append(vals, fmt.Sprintf(`"%s" : %s`, k, v.Value.(string)))
					} else if nestedMedia { // media in repeated field
						if val, ok := v.Value.(string); ok {
							splits := strings.Split(val, ".")
							end := strings.Join(splits[len(splits)-2:len(splits)], ".")
							vals = append(vals, fmt.Sprintf(`"%s" : %s.%s`, k, _var, end))
							start := strings.Join(splits[:len(splits)-2], ".")
							_pv = start
						}
					} else {
						if val, ok := v.Value.(string); ok {
							splitAt := strings.LastIndex(val, ".")
							vals = append(vals, fmt.Sprintf(`"%s" : %s.%s`, k, _var, val[splitAt+1:len(val)]))
							_pv = val[:splitAt]
						}
					}
				}
			}
		}
		return vals, _pv, nil
	}
*/
func ReturnAQLObjectBuilder(_var string, isArray bool, isParentArray bool, _map map[string]*shared.FieldDetails) ([]string, error) {
	var vals []string
	for k, v := range _map {
		if v.Value != nil {
			switch v.FieldType {
			case "repeated":
				if !utility.ArrayContains([]string{"data", "meta", "created_by", "last_modified_by"}, k) { // skip for data object
					_nestedVar := utility.RandomVariableGenerator(4)
					_returns, err := ReturnAQLObjectBuilder(_nestedVar, true, isArray, v.Value.(map[string]*shared.FieldDetails))
					if err != nil {
						return nil, err
					}
					_array := fmt.Sprintf("%s.`%s`", _var, k)
					q := fmt.Sprintf(`"%s" : ( FOR %s in NOT_NULL(%s) ? %s : [] RETURN { %s } )`, k, _nestedVar, _array, _array, strings.Join(_returns, ", "))
					vals = append(vals, q)
				} else {
					_returns, err := ReturnAQLObjectBuilder(fmt.Sprintf(`%s.%s`, _var, k), false, false, v.Value.(map[string]*shared.FieldDetails))
					if err != nil {
						return nil, err
					}
					q := fmt.Sprintf(`"%s" : { %s }`, k, strings.Join(_returns, ", "))
					vals = append(vals, q)
				}
				break
			case "object":
				if !utility.ArrayContains([]string{"data", "meta", "created_by", "last_modified_by"}, k) { // skip for data object
					_nestedVar := fmt.Sprintf("%s.`%s`", _var, k)
					_returns, err := ReturnAQLObjectBuilder(_nestedVar, true, isArray, v.Value.(map[string]*shared.FieldDetails))
					if err != nil {
						return nil, err
					}
					q := fmt.Sprintf(`"%s" : { %s } `, k, strings.Join(_returns, ", "))
					vals = append(vals, q)
				} else {
					_returns, err := ReturnAQLObjectBuilder(fmt.Sprintf(`%s.%s`, _var, k), false, false, v.Value.(map[string]*shared.FieldDetails))
					if err != nil {
						return nil, err
					}
					q := fmt.Sprintf(`"%s" : { %s }`, k, strings.Join(_returns, ", "))
					vals = append(vals, q)
				}
				break
			case "media", "geo":
				if v.Validation != nil && v.Validation.IsGallery { // multiple media is an array
					_nestedVar := utility.RandomVariableGenerator(4)
					_returns, err := ReturnAQLObjectBuilder(_nestedVar, true, isArray, v.Value.(map[string]*shared.FieldDetails))
					if err != nil {
						return nil, err
					}
					_array := fmt.Sprintf("%s.`%s`", _var, k)
					q := fmt.Sprintf(`"%s" : ( FOR %s in NOT_NULL(%s) ? %s : [] RETURN { %s } )`, k, _nestedVar, _array, _array, strings.Join(_returns, ", "))
					vals = append(vals, q)
				} else {
					var newVar string
					if isArray {
						newVar = _var
					} else {
						newVar = fmt.Sprintf("%s.`%s`", _var, k)
					}
					_returns, err := ReturnAQLObjectBuilder(newVar, false, isArray, v.Value.(map[string]*shared.FieldDetails))
					if err != nil {
						return nil, err
					}
					vals = append(vals, fmt.Sprintf(`"%s" : { %s }`, k, strings.Join(_returns, ", ")))
				}
				break
			case "multiline":
				var newVar string
				if isArray {
					newVar = _var
				} else {
					newVar = fmt.Sprintf("%s.`%s`", _var, k)
				}
				_returns, err := ReturnAQLObjectBuilder(newVar, false, isArray, v.Value.(map[string]*shared.FieldDetails))
				if err != nil {
					return nil, err
				}
				vals = append(vals, fmt.Sprintf(`"%s" : { %s }`, k, strings.Join(_returns, ", ")))
				break
			default:
				if val, ok := v.Value.(string); ok {
					var value string
					varSplit := strings.Split(val, ".")
					if isArray { // its not nested
						//_varSplit := strings.Split(_var, ".")
						end := fmt.Sprintf("`%s`", strings.Join(varSplit[len(varSplit)-1:len(varSplit)], "`.`"))
						//start := strings.Join(_varSplit[:len(_varSplit)-1], ".")
						value = fmt.Sprintf("%s.%s", _var, end)
					} else if !isArray && isParentArray { // its not nested
						//_varSplit := strings.Split(_var, ".")
						end := fmt.Sprintf("`%s`", strings.Join(varSplit[len(varSplit)-2:len(varSplit)], "`.`"))
						//start := strings.Join(_varSplit[:len(_varSplit)-1], ".")
						value = fmt.Sprintf("%s.%s", _var, end)
					} else {
						end := fmt.Sprintf("`%s`", strings.Join(varSplit[len(varSplit)-1:len(varSplit)], "`.`"))
						start := fmt.Sprintf("`%s`", strings.Join(varSplit[:len(varSplit)-1], "`.`"))
						value = fmt.Sprintf("%s.%s", start, end)
					}
					q := fmt.Sprintf(`"%s" : %s`, k, value)
					vals = append(vals, q)
				}
			}
		}
	}
	if len(vals) > 0 {
		// inject id field for gqlgen support
		vals = append(vals, fmt.Sprintf(`"id" : %s.id`, strings.TrimSuffix(_var, ".data")))
	}
	return vals, nil
}
