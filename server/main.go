package main

import "net/http"

func main() {
	http.HandleFunc("/", fetchCurrency)

	http.ListenAndServe(":8080", nil)
}

func fetchCurrency(w http.ResponseWriter, req *http.Request) {

}
