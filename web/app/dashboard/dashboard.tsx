import type React from "react";
import { useState, useMemo, useEffect, Suspense } from "react";
import { addMonth, addYear, dayEnd } from "@formkit/tempo";

import { StockChart, type ChartData } from "./stockChart";
import { Pagination } from "./pagination";
import { StockCard } from "./stockCard";

import {
  useCreateNotificationMutation,
  useGetNotificationQuery,
  useGetSymbolsSuspenseQuery,
} from "../gen/graphql";
import { PeriodSelector, type Period } from "./periodSelector";
import { NotificationSection } from "./notification";
import type { ApolloError } from "@apollo/client";

// 1ページあたりの表示件数
const ITEMS_PER_PAGE = 4;

const periodToTitle = (period: Period) => {
  switch (period) {
    case "1M":
      return "過去1ヶ月の価格推移";
    case "6M":
      return "過去6ヶ月の価格推移";
    case "1Y":
      return "過去12ヶ月の価格推移";
    case "5Y":
      return "過去5年の価格推移";
    default:
      period satisfies never;
      return "";
  }
};

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

const useGetSymbols = () => {
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

const isError = (error?: ApolloError) => {
  if (error === undefined) return false;
  const extension = error.cause?.extensions;
  const extensions = Array.isArray(extension) ? extension : extension == null ? [] : [extension];
  return extensions.find((e) => e.code) ? true : false;
};
const useGetNotification = () => {
  const { data, loading, error } = useGetNotificationQuery();
  return {
    data,
    loading,
    unAuthorization: isError(error),
  };
};

const useCreateNotification = () => {
  const parseTime = (time: string) => {
    const hour = Number(time.slice(0, 2));
    const date = new Date();
    date.setHours(hour);
    return date;
  };
  const [mutate] = useCreateNotificationMutation();
  const handleCreateNotification = (notification: {
    id: number;
    time: string;
    selectedStockSymbols: string[];
  }) => {
    mutate({
      variables: {
        createNotificationInput: {
          symbols: notification.selectedStockSymbols,
          time: parseTime(notification.time).toISOString(),
        },
      },
    }).then((r) => {
      console.warn(r.data);
    });
  };
  return { handleCreateNotification };
};

/**
 * ダッシュボード全体のページコンポーネント
 */
function DashboardPage() {
  const {
    symbols,
    selectedSymbol,
    setSelectedSymbol,
    selectedPeriod,
    setSelectedPeriod,
    chartData,
  } = useGetSymbols();
  const [searchQuery, setSearchQuery] = useState(""); // 検索クエリ
  const [currentPage, setCurrentPage] = useState(1); // 現在のページ番号

  const { data, loading, unAuthorization } = useGetNotification();
  const { handleCreateNotification } = useCreateNotification();

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
  const handleCardClick = (symbol: { symbol: string; shortName: string }) => {
    setSelectedSymbol(symbol);
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

  const title = useMemo(() => {
    if (selectedSymbol) {
      return `${selectedSymbol.shortName} - ${periodToTitle(selectedPeriod)}`;
    }
    return "銘柄を選択してください";
  }, [selectedSymbol, selectedPeriod]);

  return (
    <Suspense fallback={<div className="text-center">Loading...</div>}>
      <div className="p-4 sm:p-8 bg-gray-100 dark:bg-gray-900 min-h-screen font-sans lg:max-w-7xl mx-auto">
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
                  isSelected={stock.symbol === selectedSymbol?.symbol}
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
        <div className="mt-8 mb-8">
          {/* 期間選択 */}
          <PeriodSelector
            currentPeriod={selectedPeriod}
            onPeriodChange={(period: Period) => setSelectedPeriod(period)}
          />
          {/* チャートの上にマージンを追加 */}
          <StockChart data={chartData.data} title={title} tickFormatter={chartData.formatter} />
        </div>

        {import.meta.env.VITE_ENABLE_STOCK_NOTIFICATION === "true" && (
          <NotificationSection
            user={{ displayName: "ゲストユーザー" }}
            isAuthorized={!unAuthorization}
            handleGoogleLogin={() => {}}
            handleLogout={() => {}}
            handleAddNotification={handleCreateNotification}
            handleDeleteNotification={() => {}}
            allStocks={symbols?.map(({ symbol, shortName }) => ({ symbol, name: shortName })) || []}
            notifications={[]}
          />
        )}

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
    <div className="bg-gray-100 dark:bg-gray-900">
      <DashboardPage />
    </div>
  );
}
