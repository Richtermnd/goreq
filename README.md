# Goreq


Lib to parse requests data from different sources like path, query and body (while only json)


### !!! It's not final version.

Example struct and handler

``` go
type Request struct {
    Id        int    `path:"id"`
    FirstName string `query:"firstName"`
    LastName  string `json: "lastName"`
}

http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := goreq.Decode(r, &req)
		if err != nil {
			goreq.SendError(w, err)
			return
		}
		fmt.Printf("%+v\n", req)  // see result
        w.WriteHeader(http.StatusOK)
	})
```
