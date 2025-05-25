package notifystock

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader/v7"
)

func Ptr[V any](v V) *V {
	return &v
}

type DataLoader struct {
	SymbolDetail *dataloader.Loader[string, *SymbolDetail]
}

func NewDataLoader(
	symbolRepository *SymbolRepository,
) *DataLoader {
	symbolDetail := dataloader.NewBatchedLoader(func(ctx context.Context, keys []string) []*dataloader.Result[*SymbolDetail] {
		symbols, err := symbolRepository.GetBySymbols(ctx, keys)
		if err != nil {
			return []*dataloader.Result[*SymbolDetail]{
				{Data: nil, Error: err},
			}
		}
		m := make(map[string]*SymbolDetail, len(symbols))
		for _, symbol := range symbols {
			m[symbol.Symbol] = &symbol
		}

		results := make([]*dataloader.Result[*SymbolDetail], len(keys))
		for i, key := range keys {
			if symbol, ok := m[key]; ok {
				results[i] = &dataloader.Result[*SymbolDetail]{
					Data:  symbol,
					Error: nil,
				}
			} else {
				results[i] = &dataloader.Result[*SymbolDetail]{
					Data:  nil,
					Error: fmt.Errorf("symbol %s not found", key),
				}
			}
		}
		return results
	}, dataloader.WithCache(Ptr(dataloader.NoCache[string, *SymbolDetail]{})))
	return &DataLoader{
		SymbolDetail: symbolDetail,
	}
}
