package goreq_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/Richtermnd/goreq"
)

func TestDecodeJson(t *testing.T) {
	type JsonRequest struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	testCases := []struct {
		desc     string
		path     map[string]interface{}
		query    map[string]interface{}
		body     map[string]interface{}
		expected JsonRequest
		wantErr  bool
	}{
		{
			desc: "simple json",
			body: map[string]interface{}{
				"id":         1,
				"first_name": "John",
				"last_name":  "Doe",
			},
			expected: JsonRequest{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			desc:  "Wrong data source",
			path:  map[string]interface{}{"id": 1},
			query: map[string]interface{}{"first_name": "John"},
			body:  map[string]interface{}{"last_name": "Doe"},
			expected: JsonRequest{
				LastName: "Doe",
			},
		},
		{
			desc:     "Wrong data type",
			body:     map[string]interface{}{"id": "1", "first_name": false, "last_name": 1},
			expected: JsonRequest{},
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := buildRequest(tc.path, tc.query, tc.body)

			var actual JsonRequest
			err := goreq.Decode(req, &actual)
			if err != nil {
				if tc.wantErr {
					return
				}
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}

func TestDecodePath(t *testing.T) {
	type PathRequest struct {
		Id        int    `path:"id"`
		FirstName string `path:"first_name"`
		LastName  string `path:"last_name"`
	}
	testCases := []struct {
		desc     string
		path     map[string]interface{}
		query    map[string]interface{}
		body     map[string]interface{}
		expected PathRequest
		wantErr  bool
	}{
		{
			desc: "Valid request",
			path: map[string]interface{}{
				"id":         1,
				"first_name": "John",
				"last_name":  "Doe",
			},
			expected: PathRequest{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			desc:  "Wrong data source",
			path:  map[string]interface{}{"id": 1},
			query: map[string]interface{}{"first_name": "John"},
			body:  map[string]interface{}{"last_name": "Doe"},
			expected: PathRequest{
				Id: 1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := buildRequest(tc.path, tc.query, tc.body)

			var actual PathRequest
			err := goreq.Decode(req, &actual)
			if err != nil {
				if tc.wantErr {
					return
				}
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}

func TestDecodeQuery(t *testing.T) {
	type QueryRequest struct {
		Id        int    `query:"id"`
		FirstName string `query:"first_name"`
		LastName  string `query:"last_name"`
	}
	testCases := []struct {
		desc     string
		path     map[string]interface{}
		query    map[string]interface{}
		body     map[string]interface{}
		expected QueryRequest
		wantErr  bool
	}{
		{
			desc: "Valid request",
			query: map[string]interface{}{
				"id":         1,
				"first_name": "John",
				"last_name":  "Doe",
			},
			expected: QueryRequest{
				Id:        1,
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			desc:  "Wrong data source",
			path:  map[string]interface{}{"id": 1},
			query: map[string]interface{}{"first_name": "John"},
			body:  map[string]interface{}{"last_name": "Doe"},
			expected: QueryRequest{
				FirstName: "John",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := buildRequest(tc.path, tc.query, tc.body)

			var actual QueryRequest
			err := goreq.Decode(req, &actual)
			if err != nil {
				if tc.wantErr {
					return
				}
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}

func buildRequest(path, query, body map[string]interface{}) *http.Request {
	var reqBody io.Reader
	if body != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		reqBody = buf
	}
	req, err := http.NewRequest(http.MethodPost, "http://example.com", reqBody)
	if err != nil {
		panic(err)
	}
	addPathParams(req, path)
	addQueryParams(req, query)
	return req
}

func addPathParams(req *http.Request, path map[string]interface{}) {
	if path == nil {
		return
	}
	for k, v := range path {
		req.SetPathValue(k, fmt.Sprintf("%v", v))
	}
}

func addQueryParams(req *http.Request, query map[string]interface{}) {
	if query == nil {
		return
	}
	for k, v := range query {
		newQuery := req.URL.Query()
		newQuery.Set(k, fmt.Sprintf("%v", v))
		req.URL.RawQuery = newQuery.Encode()
	}
}
