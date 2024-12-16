package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	Client http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client: http.Client{},
	}
}

func (c *HTTPClient) FetchCurrentStock() error {
	now := time.Now().Unix()
	URL := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/^N225?period1=%d&period2=%d&region=US", now, now)
	fmt.Println(URL)

	res, err := c.Client.Get(URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
