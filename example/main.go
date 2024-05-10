package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Richtermnd/goreq"
)

type Request struct {
	Id        int    `path:"id"`
	FirstName string `query:"firstName"`
	LastName  string `json:"lastName"`
}

func startServer() {
	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		// decode
		var req Request
		err := goreq.Decode(r, &req)
		if err != nil {
			goreq.SendError(w, err)
			return
		}
		fmt.Printf("%+v\n", req)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func main() {
	go startServer()

	// build request
	mapBody := map[string]interface{}{
		"id":        999,   // trash data to demonstrate that goreq has source priority.
		"firstName": "xdd", // trash data to demonstrate that goreq has source priority.
		"lastName":  "Doe",
	}
	data, _ := json.Marshal(mapBody)
	req, _ := http.NewRequest("POST", "http://localhost:8080/1?firstName=John", bytes.NewReader(data))

	// send request
	http.DefaultClient.Do(req)

	// build bad request
	mapBody = map[string]interface{}{
		"lastName": 123, // wrong type
	}
	data, _ = json.Marshal(mapBody)
	req, _ = http.NewRequest("POST", "http://localhost:8080/NotInt?firstName=John", bytes.NewReader(data))

	// send request
	http.DefaultClient.Do(req)

	// wait interrupt
	ch := make(chan struct{})
	<-ch

}
