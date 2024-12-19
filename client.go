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

type ChartResponse struct {
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
	PriceHint            int                  `json:"priceHint"`
	CurrentTradingPeriod CurrentTradingPeriod `json:"currentTradingPeriod"`
	DataGranularity      string               `json:"dataGranularity"`
	Range                string               `json:"range"`
	ValidRanges          []string             `json:"validRanges"`
}
type Quote struct {
	Open   []float64 `json:"open"`
	Close  []float64 `json:"close"`
	Volume []int     `json:"volume"`
	High   []float64 `json:"high"`
	Low    []float64 `json:"low"`
}
type Adjclose struct {
	Adjclose []float64 `json:"adjclose"`
}
type Indicators struct {
	Quote    []Quote    `json:"quote"`
	Adjclose []Adjclose `json:"adjclose"`
}
type Result struct {
	Meta       Meta       `json:"meta"`
	Timestamp  []int      `json:"timestamp"`
	Indicators Indicators `json:"indicators"`
}
type Chart struct {
	Result []Result `json:"result"`
	Error  any      `json:"error"`
}

func (c *FinanceClient) FetchCurrentStock(symbol string) (*ChartResponse, error) {
	now := time.Now()
	return c.FetchStock(symbol, now, now)
}

func (c *FinanceClient) FetchStock(symbol string, beggingOfPeriod, endOfPeriod time.Time) (*ChartResponse, error) {
	URL := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&events=div|split|earn&interval=1d&region=US",
		symbol, beggingOfPeriod.Unix(), endOfPeriod.Unix())
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
	fmt.Println(string(body))
	var chart ChartResponse
	if err := json.Unmarshal(body, &chart); err != nil {
		return nil, err
	}
	return &chart, nil
}
