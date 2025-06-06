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
		sym := make([]*SymbolDetail, 0, len(symbols))
		for i := range symbols {
			sym = append(sym, &symbols[i])
		}
		return createResult(sym, keys)
	}, dataloader.WithCache(Ptr(dataloader.NoCache[string, *SymbolDetail]{})))
	return &DataLoader{
		SymbolDetail: symbolDetail,
	}
}

type DataLoaderKey interface {
	Key() string
}

func createResult[T DataLoaderKey](data []T, keys []string) []*dataloader.Result[T] {
	m := make(map[string]T, len(data))
	for _, v := range data {
		m[v.Key()] = v
	}
	results := make([]*dataloader.Result[T], len(keys))
	for i, key := range keys {
		if v, ok := m[key]; ok {
			results[i] = &dataloader.Result[T]{Data: v, Error: nil}
		} else {
			results[i] = &dataloader.Result[T]{Error: fmt.Errorf("key %v not found", key)}
		}
	}
	return results
}
