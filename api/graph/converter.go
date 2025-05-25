package graph

import (
	"github.com/heyjun3/notify-stock/graph/model"
	notify "github.com/heyjun3/notify-stock/internal"
)

func convertToSymbolDetail(symbol *notify.SymbolDetail) *model.SymbolDetail {
	if symbol == nil {
		return nil
	}
	return &model.SymbolDetail{
		Symbol:         symbol.Symbol,
		ShortName:      symbol.ShortName,
		LongName:       symbol.LongName,
		Price:          symbol.MarketPrice.InexactFloat64(),
		Change:         symbol.Change(),
		ChangePercent:  symbol.ChangePercent(),
		CurrencySymbol: symbol.Currency.Symbol(),
	}
}
func convertToSymbolDetails(symbols []*notify.SymbolDetail) []*model.SymbolDetail {
	if len(symbols) == 0 {
		return nil
	}
	details := make([]*model.SymbolDetail, 0, len(symbols))
	for _, symbol := range symbols {
		details = append(details, convertToSymbolDetail(symbol))
	}
	return details
}
