package goreq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// Decode decode request to destination
//
// dest must be a pointer to valid struct
func Decode(r *http.Request, dest interface{}) (err error) {
	// catch panic cuz idgaf how make it better without pain in ass
	defer func() {
		if r := recover(); r != nil {
			err = newGoreqError("Bad request", http.StatusBadRequest, ErrBadRequest)
		}
	}()
	return decode(r, dest)
}

// SendError write error in response body as {"msg": "message text"}
func SendError(w http.ResponseWriter, err error) {
	goreqErr, ok := err.(GoreqError)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newGoreqError("bad request", http.StatusBadRequest, err))
		return
	}
	w.WriteHeader(goreqErr.HttpCode)
	json.NewEncoder(w).Encode(goreqErr)
}

func decode(r *http.Request, dest interface{}) error {
	if !isPointerToStruct(dest) {
		return ErrNotPointerToStruct
	}
	// get struct schema
	schema := analyzeDestination(dest)

	// try to fill dest with values from request
	if err := fillDestination(r, schema, dest); err != nil {
		return err
	}

	return nil
}

func fillDestination(req *http.Request, schema structSchema, dest interface{}) error {
	elem := reflect.ValueOf(dest).Elem()

	// Get and set path value
	for _, fieldSchema := range schema.path {
		// get and check field
		field := elem.FieldByName(fieldSchema.name)
		if err := isValidField(field); err != nil {
			return err
		}

		// get and set path value
		pathValue := req.PathValue(fieldSchema.sourceName)
		value, err := parseByKind(fieldSchema.kind, pathValue)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(value).Convert(field.Type()))
	}

	// Get and set query value
	for _, fieldSchema := range schema.query {
		// get and check field
		field := elem.FieldByName(fieldSchema.name)
		if err := isValidField(field); err != nil {
			return err
		}

		// get and set query value
		queryValue := req.URL.Query().Get(fieldSchema.sourceName)
		if queryValue == "" {
			continue
		}
		value, err := parseByKind(fieldSchema.kind, queryValue)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(value).Convert(field.Type()))
	}

	// Check body
	if req.Body == nil {
		return nil
	}

	// Get and set json values
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(req.Body).Decode(&jsonBody)
	if err != nil {
		return err
	}

	for _, fieldSchema := range schema.json {
		field := elem.FieldByName(fieldSchema.name)
		if err := isValidField(field); err != nil {
			return err
		}
		jsonField, ok := jsonBody[fieldSchema.sourceName]
		if !ok {
			continue
		}
		field.Set(reflect.ValueOf(jsonField).Convert(field.Type()))
	}

	return nil
}

// isPointerToStruct check if dest is a pointer to struct
func isPointerToStruct(dest interface{}) bool {
	return reflect.TypeOf(dest).Kind() == reflect.Ptr && reflect.TypeOf(dest).Elem().Kind() == reflect.Struct
}

// isValidField check if field is valid and not private
func isValidField(field reflect.Value) error {
	if !field.IsValid() {
		return fmt.Errorf("%w: invalid field", ErrInvalidDestination)
	}
	if !field.CanSet() {
		return fmt.Errorf("%w: cannot set field", ErrInvalidDestination)
	}
	return nil
}
