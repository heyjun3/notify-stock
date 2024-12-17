package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClientInterface interface {
	Get(string) (*http.Response, error)
}

type FinanceClient struct {
	Client HTTPClientInterface
}

func NewFinanceClient(client HTTPClientInterface) *FinanceClient {
	return &FinanceClient{
		Client: client,
	}
}

type CurrentStokResponse struct {
	Chart Chart `json:"chart"`
}
type Pre struct {
	Timezone  string `json:"timezone"`
	End       int    `json:"end"`
	Start     int    `json:"start"`
	Gmtoffset int    `json:"gmtoffset"`
}
type Regular struct {
	Timezone  string `json:"timezone"`
	End       int    `json:"end"`
	Start     int    `json:"start"`
	Gmtoffset int    `json:"gmtoffset"`
}
type Post struct {
	Timezone  string `json:"timezone"`
	End       int    `json:"end"`
	Start     int    `json:"start"`
	Gmtoffset int    `json:"gmtoffset"`
}
type CurrentTradingPeriod struct {
	Pre     Pre     `json:"pre"`
	Regular Regular `json:"regular"`
	Post    Post    `json:"post"`
}
type Meta struct {
	Currency             string               `json:"currency"`
	Symbol               string               `json:"symbol"`
	ExchangeName         string               `json:"exchangeName"`
	FullExchangeName     string               `json:"fullExchangeName"`
	InstrumentType       string               `json:"instrumentType"`
	FirstTradeDate       int                  `json:"firstTradeDate"`
	RegularMarketTime    int                  `json:"regularMarketTime"`
	HasPrePostMarketData bool                 `json:"hasPrePostMarketData"`
	Gmtoffset            int                  `json:"gmtoffset"`
	Timezone             string               `json:"timezone"`
	ExchangeTimezoneName string               `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64              `json:"regularMarketPrice"`
	FiftyTwoWeekHigh     float64              `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow      float64              `json:"fiftyTwoWeekLow"`
	RegularMarketDayHigh float64              `json:"regularMarketDayHigh"`
	RegularMarketDayLow  float64              `json:"regularMarketDayLow"`
	RegularMarketVolume  int                  `json:"regularMarketVolume"`
	LongName             string               `json:"longName"`
	ShortName            string               `json:"shortName"`
	ChartPreviousClose   float64              `json:"chartPreviousClose"`
	PreviousClose        float64              `json:"previousClose"`
	Scale                int                  `json:"scale"`
	PriceHint            int                  `json:"priceHint"`
	CurrentTradingPeriod CurrentTradingPeriod `json:"currentTradingPeriod"`
	DataGranularity      string               `json:"dataGranularity"`
	Range                string               `json:"range"`
	ValidRanges          []string             `json:"validRanges"`
}
type Quote struct {
}
type Indicators struct {
	Quote []Quote `json:"quote"`
}
type Result struct {
	Meta       Meta       `json:"meta"`
	Indicators Indicators `json:"indicators"`
}
type Chart struct {
	Result []Result `json:"result"`
	Error  any      `json:"error"`
}

func (c *FinanceClient) FetchCurrentStock(symbol string) (*CurrentStokResponse, error) {
	now := time.Now().Unix()
	URL := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&region=US",
		symbol, now, now)
	fmt.Println(URL)

	res, err := c.Client.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var current CurrentStokResponse
	if err := json.Unmarshal(body, &current); err != nil {
		return nil, err
	}
	fmt.Println(current.Chart.Result[0].Meta.RegularMarketPrice)
	return &current, nil
}
