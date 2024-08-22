package sql

import (
	"context"
	"fmt"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	"github.com/graph-gophers/dataloader"
)

type FieldClassification struct {
	MultilineFields []string
	DoubleFields    []string
	PictureField    []string
	GalleryField    []string
	ListFields      []string
	RepeatedFields  map[string][]*protobuff.FieldInfo
}

func (S *ProjectSqlDriver) RelationshipDataLoaderBytes(ctx context.Context, param *shared.CommonSystemParams, connection map[string]interface{}) ([]byte, error) {
	// query relations and find all docs
	query, _, err := BuildCombinedRelationQuery("--removed", "--removed", param)
	if err != nil {
		return nil, err
	}

	var queryResults []map[string]interface{}
	err = S.ORM.NewRaw(*query).Scan(ctx, &queryResults)
	if err != nil {
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

	finalResults := make(map[string][]*shared.DefaultDocumentStructure)
	// format query results to finalResults
	for _, res := range queryResults {
		if val, ok := res["sys_key"].([]byte); ok {
			key := string(val)
			if docs, ok := finalResults[key]; ok {
				doc, err := CommonDocTransformation(param.Model, "en", res, &classification)
				if err != nil {
					fmt.Println(err.Error())
				}
				docs = append(docs, doc)
				finalResults[key] = docs
			} else {
				doc, err := CommonDocTransformation(param.Model, "en", res, &classification)
				if err != nil {
					fmt.Println(err.Error())
				}
				finalResults[key] = []*shared.DefaultDocumentStructure{doc}
			}
		}
	}

	keys := param.DocumentIDs

	var results []*dataloader.Result
	switch connection["relation_type"] {
	case "has_many":
		// prepare the result
		for _, id := range keys {
			result := dataloader.Result{
				Data:  finalResults[id],
				Error: nil,
			}
			results = append(results, &result)
		}
		break
	case "has_one":
		// prepare the result
		for _, id := range keys {
			if len(finalResults[id]) > 0 {
				results = append(results, &dataloader.Result{
					Data:  finalResults[id][0], // because it has only one
					Error: nil,
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  nil, // because it has only one
					Error: nil,
				})
			}
		}
		break
	}

	return []byte{}, nil
}

func (S *ProjectSqlDriver) RelationshipDataLoader(ctx context.Context, param *shared.CommonSystemParams, connection map[string]interface{}) (interface{}, error) {
	// query relations and find all docs
	query, _, err := BuildCombinedRelationQuery("--removed", "--removed", param)
	if err != nil {
		return nil, err
	}

	queryResults := []map[string]interface{}{}
	err = S.ORM.NewRaw(*query).Scan(ctx, &queryResults)
	if err != nil {
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

	finalResults := make(map[string][]*shared.DefaultDocumentStructure)
	// format query results to finalResults
	for _, res := range queryResults {
		if val, ok := res["sys_key"].([]byte); ok {
			key := string(val)
			if docs, ok := finalResults[key]; ok {
				doc, err := CommonDocTransformation(param.Model, "en", res, &classification)
				if err != nil {
					fmt.Println(err.Error())
				}
				docs = append(docs, doc)
				finalResults[key] = docs
			} else {
				doc, err := CommonDocTransformation(param.Model, "en", res, &classification)
				if err != nil {
					fmt.Println(err.Error())
				}
				finalResults[key] = []*shared.DefaultDocumentStructure{doc}
			}
		}
	}

	keys := param.DocumentIDs

	var results []*dataloader.Result
	switch connection["relation_type"] {
	case "has_many":
		// prepare the result
		for _, id := range keys {
			result := dataloader.Result{
				Data:  finalResults[id],
				Error: nil,
			}
			results = append(results, &result)
		}
		break
	case "has_one":
		// prepare the result
		for _, id := range keys {
			if len(finalResults[id]) > 0 {
				results = append(results, &dataloader.Result{
					Data:  finalResults[id][0], // because it has only one
					Error: nil,
				})
			} else {
				results = append(results, &dataloader.Result{
					Data:  nil, // because it has only one
					Error: nil,
				})
			}
		}
		break
	}

	return results, nil
}
