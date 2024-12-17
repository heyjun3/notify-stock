package main

import (
	"net/http"
)

func main() {
	client := NewFinanceClient(&http.Client{})
	client.FetchCurrentStock("^N225")
}
