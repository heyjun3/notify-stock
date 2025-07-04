import { addMonth, addYear, dayEnd } from "@formkit/tempo";
import { useState, useEffect } from "react";

import { useGetSymbolsSuspenseQuery } from "~/gen/graphql";
import type { Period } from "../periodSelector";
import type { ChartData } from "../stockChart";

const periodToStart = (period: Period) => {
  switch (period) {
    case "1M":
      return addMonth(new Date(), -1);
    case "6M":
      return addMonth(new Date(), -6);
    case "1Y":
      return addYear(new Date(), -1);
    case "5Y":
      return addYear(new Date(), -5);
    default:
      period satisfies never;
      return new Date();
  }
};

type Chart = {
  data: ChartData[];
  formatter?: (v: any) => string;
};

export const useGetSymbols = () => {
  const [selectedSymbol, setSelectedSymbol] = useState<{
    symbol: string;
    shortName: string;
  } | null>(null); // 初期選択
  const [selectedPeriod, setSelectedPeriod] = useState<Period>("1Y"); // 選択された期間
  const [chartData, setChartData] = useState<Chart>({ data: [] });
  const { data } = useGetSymbolsSuspenseQuery({
    variables: {
      chartInput: {
        start: dayEnd(periodToStart(selectedPeriod)).toISOString(),
        end: dayEnd(new Date()).toISOString(),
        symbol: selectedSymbol?.symbol,
      },
    },
  });
  const symbols = data?.symbols.map((symbol) => symbol.detail);

  useEffect(() => {
    if (data && data.symbols.length > 0) {
      setSelectedSymbol(
        (symbol) =>
          symbol ?? {
            symbol: data.symbols[0].detail.symbol,
            shortName: data.symbols[0].detail.shortName,
          },
      );
    }
    if (!data) return;
    for (const symbol of data.symbols) {
      const chart = symbol.chart;
      if (chart && chart.length) {
        const formatter = (value: any) => `${symbol.detail.currencySymbol}${value}`;
        setChartData({ data: chart, formatter });
      }
    }
    return;
  }, [data]);
  return {
    symbols,
    selectedSymbol,
    setSelectedSymbol,
    selectedPeriod,
    setSelectedPeriod,
    chartData,
  };
};
