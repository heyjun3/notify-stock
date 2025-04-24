package notifystock

//go:generate enumer -type=Currency
type Currency int

const (
	_ Currency = iota
	JPY
	USD
)
