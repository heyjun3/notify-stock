package main

import (
	"net/http"
	"time"
)

func main() {
	client := NewFinanceClient(&http.Client{})
	client.FetchStock("^N225", time.Date(2024, 12, 13, 0, 0, 0, 0, time.UTC), time.Now())
}
