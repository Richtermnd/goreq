package goreq

import "reflect"

// structSchema store parsed fields of struct sorted by tags
type structSchema struct {
	path  []fieldSchema
	query []fieldSchema
	json  []fieldSchema
}

// fieldSchema store parsed field schema
type fieldSchema struct {
	name       string
	sourceName string
	kind       reflect.Kind
}

// analyzeDestination analyze destination struct
func analyzeDestination(dest interface{}) structSchema {
	schema := structSchema{
		path:  make([]fieldSchema, 0),
		query: make([]fieldSchema, 0),
		json:  make([]fieldSchema, 0),
	}
	fields := getFields(dest)
	for _, field := range fields {
		// if path tag defined
		pathTag := field.Tag.Get("path")
		if pathTag != "" && pathTag != "-" {
			schema.path = append(schema.path, fieldSchema{
				name:       field.Name,
				sourceName: pathTag,
				kind:       field.Type.Kind(),
			})
			continue
		}

		// if query tag defined
		queryTag := field.Tag.Get("query")
		if queryTag != "" && queryTag != "-" {
			schema.query = append(schema.query, fieldSchema{
				name:       field.Name,
				sourceName: queryTag,
				kind:       field.Type.Kind(),
			})
			continue
		}

		// if json tag defined
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			schema.json = append(schema.json, fieldSchema{
				name:       field.Name,
				sourceName: jsonTag,
				kind:       field.Type.Kind(),
			})
			continue
		}
	}
	return schema
}

// getFields get all struct fields
func getFields(dest interface{}) []reflect.StructField {
	elem := reflect.TypeOf(dest).Elem()
	fields := make([]reflect.StructField, 0, elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		fields = append(fields, elem.Field(i))
	}
	return fields
}
