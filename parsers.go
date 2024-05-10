package goreq

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func parseByKind(kind reflect.Kind, s string) (interface{}, error) {
	switch kind {
	case reflect.String:
		return s, nil
	case reflect.Int:
		return parseInt(s)
	case reflect.Bool:
		return parseBool(s)
	case reflect.Float32, reflect.Float64:
		return parseFloat(s)
	default:
		msg := fmt.Sprintf("can't parse %s", s)
		return nil, newGoreqError(msg, http.StatusBadRequest, ErrBadRequest)
	}
}

func parseInt(s string) (int64, error) {
	intValue, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("failed to convert %s to int", s)
		return 0, newGoreqError(msg, http.StatusBadRequest, ErrBadRequest)
	}
	return intValue, nil
}

func parseBool(s string) (bool, error) {
	boolValue, err := strconv.ParseBool(s)
	if err != nil {
		msg := fmt.Sprintf("failed to convert %s to bool", s)
		return false, newGoreqError(msg, http.StatusBadRequest, ErrBadRequest)
	}
	return boolValue, nil
}

func parseFloat(s string) (float64, error) {
	floatValue, err := strconv.ParseFloat(s, 64)
	if err != nil {
		msg := fmt.Sprintf("failed to convert %s to float", s)
		return 0, newGoreqError(msg, http.StatusBadRequest, ErrBadRequest)
	}
	return floatValue, nil
}
