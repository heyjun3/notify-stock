package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := NewFinanceClient(&http.Client{})
	res, _ := client.FetchStock("^N225", time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC), time.Now(), WithInterval("1d"))
	fmt.Println(CalcAVG(res.Chart.Result[0].Indicators.Quote[0].Close))
}
