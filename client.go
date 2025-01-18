package notifystock

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
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

type Option func(URL *url.URL) *url.URL

func WithInterval(interval string) Option {
	return func(URL *url.URL) *url.URL {
		query := URL.Query()
		query.Add("interval", interval)
		URL.RawQuery = query.Encode()
		return URL
	}
}

func (c *FinanceClient) FetchCurrentStock(symbol Symbol) (*ChartResponse, error) {
	now := time.Now()
	return c.FetchStock(symbol, now, now)
}

func (c *FinanceClient) FetchStock(symbol Symbol, beggingOfPeriod, endOfPeriod time.Time, opts ...Option) (*ChartResponse, error) {
	s, err := symbol.ForFinance()
	if err != nil {
		return nil, err
	}
	URL, err := url.Parse(fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s", s))
	if err != nil {
		return nil, err
	}
	query := URL.Query()
	query.Add("period1", strconv.Itoa(int(beggingOfPeriod.Unix())))
	query.Add("period2", strconv.Itoa(int(endOfPeriod.Unix())))
	query.Add("region", "US")

	URL.RawQuery = query.Encode()
	for _, opt := range opts {
		URL = opt(URL)
	}

	logger.Info("request", "url", URL.String())
	res, err := c.Client.Get(URL.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var chart ChartResponse
	if err := json.Unmarshal(body, &chart); err != nil {
		return nil, err
	}
	return &chart, nil
}

func CalcAVG[T cmp.Ordered](values []T) (decimal.Decimal, error) {
	d := make([]decimal.Decimal, 0, len(values))
	for _, v := range values {
		deci, err := decimal.NewFromString(fmt.Sprintf("%v", v))
		if err != nil {
			logger.Error("failed convert string to decimal")
			return decimal.Decimal{}, err
		}
		d = append(d, deci)
	}
	avg := decimal.Avg(d[0], d[1:]...)
	return avg, nil
}
