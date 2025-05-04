import type React from "react";
import { useState, useMemo, useEffect, Suspense } from "react";
import { addYear, dayEnd, format } from "@formkit/tempo";

import { StockChart } from "./stockChart";
import { Pagination } from "./pagination";
import { StockCard } from "./stockCard";

import { useGetSymbolsSuspenseQuery } from "../gen/graphql";

// 1ページあたりの表示件数
const ITEMS_PER_PAGE = 4;

const useGetSymbols = () => {
  const [selectedSymbol, setSelectedSymbol] = useState<string | null>(null); // 初期選択
  const [charts, setCharts] = useState<
    Map<
      string,
      { shortName: string; data: { name: string; price: number }[]; formatter: (v: any) => string }
    >
  >(new Map());
  const { data } = useGetSymbolsSuspenseQuery({
    variables: {
      chartInput: {
        start: dayEnd(addYear(new Date(), -1)).toISOString(),
        end: dayEnd(new Date()).toISOString(),
      },
    },
  });
  const symbols = data?.symbols.map((symbol) => symbol.detail);

  useEffect(() => {
    if (data && data.symbols.length > 0) {
      setSelectedSymbol(data.symbols[0].detail.shortName);
    }
    if (!data) return;
    const m = new Map<
      string,
      { shortName: string; data: { name: string; price: number }[]; formatter: (v: any) => string }
    >();
    for (const symbol of data.symbols) {
      const shortName = symbol.detail.shortName;
      const chart = symbol.chart
        ?.filter((data) => data != null)
        .map((data) => ({
          name: format(data.timestamp, "short"),
          price: data.close,
        }));
      const formatter = (value: any) => `${symbol.detail.currencySymbol}${value}`;
      if (chart) {
        m.set(shortName, { shortName, data: chart, formatter });
      }
    }
    setCharts(m);
    return;
  }, [data]);
  return { symbols, charts, selectedSymbol, setSelectedSymbol };
};

/**
 * ダッシュボード全体のページコンポーネント
 */
function DashboardPage() {
  const { symbols, charts, selectedSymbol, setSelectedSymbol } = useGetSymbols();
  const [searchQuery, setSearchQuery] = useState(""); // 検索クエリ
  const [currentPage, setCurrentPage] = useState(1); // 現在のページ番号

  // 検索フィルタリング
  const filteredStocks = useMemo(() => {
    if (!symbols) return [];
    const lowerCaseQuery = searchQuery.toLowerCase();
    return symbols.filter(
      (stock) =>
        stock.symbol.toLowerCase().includes(lowerCaseQuery) ||
        stock.shortName.toLowerCase().includes(lowerCaseQuery),
    );
  }, [symbols, searchQuery]);

  // ページネーション
  const totalPages = Math.ceil(filteredStocks.length / ITEMS_PER_PAGE);
  const paginatedStocks = useMemo(() => {
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    const endIndex = startIndex + ITEMS_PER_PAGE;
    return filteredStocks.slice(startIndex, endIndex);
  }, [filteredStocks, currentPage]);

  // カードクリック時のハンドラ
  const handleCardClick = (symbol: string) => {
    setSelectedSymbol(symbol);
    // ここで選択された銘柄に応じてチャートデータを更新するロジックを追加
    // 例: fetchChartData(symbol).then(data => setChartData(data));
    console.log(`Selected stock: ${symbol}`);
    // 必要であれば、チャートデータをAPIから取得し直す
    // setChartData(fetchChartDataForSymbol(symbol)); // ダミー関数
  };

  // 検索入力ハンドラ
  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value);
    setCurrentPage(1); // 検索時は1ページ目に戻す
  };

  // ページ変更ハンドラ
  const handlePageChange = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  };

  // 選択された銘柄に対応するチャートデータを取得（ダミー）
  // 実際のアプリではAPIコールなどを行う
  const { data: currentChartData, formatter } = useMemo(() => {
    // ここで selectedSymbol に基づいて適切なチャートデータを返す
    if (!selectedSymbol || !charts) return { data: [] };
    const chart = charts.get(selectedSymbol);
    return { data: chart?.data ?? [], formatter: chart?.formatter };
  }, [selectedSymbol, charts]);

  return (
    <Suspense fallback={<div className="text-center">Loading...</div>}>
      <div className="p-4 sm:p-8 bg-gray-100 dark:bg-gray-900 min-h-screen font-sans">
        <h1 className="text-2xl sm:text-3xl font-bold mb-6 text-gray-900 dark:text-white">
          株価ダッシュボード
        </h1>

        {/* 検索バー */}
        <div className="mb-6">
          <input
            type="text"
            placeholder="銘柄シンボル or 会社名で検索..."
            value={searchQuery}
            onChange={handleSearchChange}
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          />
        </div>

        {/* 株価カード一覧 */}
        {paginatedStocks.length > 0 ? (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
              {paginatedStocks.map((stock) => (
                <StockCard
                  key={stock.symbol}
                  stock={stock}
                  isSelected={stock.shortName === selectedSymbol}
                  onClick={handleCardClick}
                />
              ))}
            </div>
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={handlePageChange}
            />
          </>
        ) : (
          <div className="text-center text-gray-500 dark:text-gray-400 py-10">
            該当する銘柄が見つかりません。
          </div>
        )}

        {/* 株価チャート */}
        <div className="mt-8">
          {" "}
          {/* チャートの上にマージンを追加 */}
          <StockChart
            data={currentChartData}
            selectedSymbol={selectedSymbol}
            tickFormatter={formatter}
          />
        </div>

        {/* フッター等、他の要素をここに追加可能 */}
        <footer className="mt-8 text-center text-gray-500 dark:text-gray-400 text-sm">
          データはダミーです。実際の取引には使用しないでください。
        </footer>
      </div>
    </Suspense>
  );
}

// --- Main App Component --- (変更なし)
export default function App() {
  return (
    <div>
      <DashboardPage />
    </div>
  );
}
