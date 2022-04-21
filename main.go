package main

import (
	"encoding/json"
	"fmt"
	"go-redis-postgres/products"
	"net/http"
)

func main() {
	http.HandleFunc("/products", httpHandler)
	http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {

	response, err := products.GetProducts()

	if err != nil {

		fmt.Fprintf(w, err.Error()+"\r\n")

	} else {

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		if err := enc.Encode(response); err != nil {
			fmt.Println(err.Error())
		}

	}
}
