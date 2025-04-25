import type React from "react";
import { useState, useMemo } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  AreaChart,
  Area,
} from "recharts";

// --- Mock Data ---
// ダミーデータを増やしてページネーションを分かりやすくする
const initialStockData = [
  {
    symbol: "AAPL",
    name: "Apple Inc.",
    price: 175.32,
    change: "+1.25",
    changePercent: "+0.72%",
    volume: "98.7M",
    marketCap: "2.75T",
  },
  {
    symbol: "GOOGL",
    name: "Alphabet Inc.",
    price: 2850.88,
    change: "-15.40",
    changePercent: "-0.54%",
    volume: "1.5M",
    marketCap: "1.90T",
  },
  {
    symbol: "MSFT",
    name: "Microsoft Corp.",
    price: 340.15,
    change: "+2.80",
    changePercent: "+0.83%",
    volume: "35.2M",
    marketCap: "2.53T",
  },
  {
    symbol: "AMZN",
    name: "Amazon.com, Inc.",
    price: 140.5,
    change: "-0.95",
    changePercent: "-0.67%",
    volume: "55.6M",
    marketCap: "1.43T",
  },
  {
    symbol: "TSLA",
    name: "Tesla, Inc.",
    price: 780.9,
    change: "+10.10",
    changePercent: "+1.31%",
    volume: "25.1M",
    marketCap: "785B",
  },
  {
    symbol: "META",
    name: "Meta Platforms, Inc.",
    price: 210.75,
    change: "-2.50",
    changePercent: "-1.17%",
    volume: "20.3M",
    marketCap: "570B",
  },
  {
    symbol: "NVDA",
    name: "NVIDIA Corporation",
    price: 235.5,
    change: "+5.20",
    changePercent: "+2.26%",
    volume: "45.8M",
    marketCap: "588B",
  },
  {
    symbol: "JPM",
    name: "JPMorgan Chase & Co.",
    price: 145.2,
    change: "+0.80",
    changePercent: "+0.55%",
    volume: "12.6M",
    marketCap: "425B",
  },
  {
    symbol: "V",
    name: "Visa Inc.",
    price: 225.4,
    change: "-1.10",
    changePercent: "-0.49%",
    volume: "8.9M",
    marketCap: "470B",
  },
  {
    symbol: "WMT",
    name: "Walmart Inc.",
    price: 142.8,
    change: "+0.30",
    changePercent: "+0.21%",
    volume: "9.5M",
    marketCap: "385B",
  },
  {
    symbol: "PG",
    name: "Procter & Gamble Co.",
    price: 155.6,
    change: "+0.90",
    changePercent: "+0.58%",
    volume: "7.2M",
    marketCap: "375B",
  },
  {
    symbol: "MA",
    name: "Mastercard Incorporated",
    price: 360.1,
    change: "-2.20",
    changePercent: "-0.61%",
    volume: "4.1M",
    marketCap: "350B",
  },
];

// ダミーのチャートデータ（選択された銘柄に応じて変化させる想定）
const chartData = [
  { name: "1月", price: 160 },
  { name: "2月", price: 165 },
  { name: "3月", price: 170 },
  { name: "4月", price: 168 },
  { name: "5月", price: 175 },
  { name: "6月", price: 180 },
  { name: "7月", price: 178 },
  { name: "8月", price: 185 },
  { name: "9月", price: 182 },
  { name: "10月", price: 175 },
  { name: "11月", price: 178 },
  { name: "12月", price: 175 },
];

// 1ページあたりの表示件数
const ITEMS_PER_PAGE = 4;

// --- Components ---

/**
 * 株価チャートを表示するコンポーネント
 * @param {object} props - コンポーネントのプロパティ
 * @param {Array<object>} props.data - チャートデータ
 * @param {string} props.selectedSymbol - 選択中の銘柄シンボル
 */
function StockChart({
  data,
  selectedSymbol,
}: { data: { name: string; price: number }[]; selectedSymbol: string | null }) {
  return (
    <div className="bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-md h-[300px] sm:h-[400px]">
      <h3 className="text-lg sm:text-xl font-semibold mb-4 text-gray-800 dark:text-gray-200">
        {selectedSymbol ? `${selectedSymbol} - 過去12ヶ月の価格推移` : "銘柄を選択してください"}
      </h3>
      {selectedSymbol && data.length > 0 ? (
        <ResponsiveContainer width="100%" height="85%">
          <AreaChart
            data={data}
            margin={{
              top: 5,
              right: 20, // 右側のスペースを確保
              left: 0, // 左側のスペースを調整
              bottom: 20, // 下部のスペースを確保
            }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke="#4B5563" /> {/* グリッド線の色を調整 */}
            <XAxis
              dataKey="name"
              tick={{ fill: "#9CA3AF" }} // 目盛りの色
              tickLine={{ stroke: "#9CA3AF" }} // 目盛り線の色
              axisLine={{ stroke: "#9CA3AF" }} // 軸線の色
              padding={{ left: 10, right: 10 }} // X軸の両端にパディングを追加
              angle={-30} // ラベルを斜めにする
              textAnchor="end" // ラベルのアンカーを調整
              height={50} // X軸の高さを確保
              interval={0} // ラベルを常に表示（データ数が多い場合は調整）
            />
            <YAxis
              tick={{ fill: "#9CA3AF" }}
              tickLine={{ stroke: "#9CA3AF" }}
              axisLine={{ stroke: "#9CA3AF" }}
              domain={["auto", "auto"]} // Y軸の範囲を自動調整
              tickFormatter={(value) => `$${value}`} // ドル記号を追加
            />
            <Tooltip
              contentStyle={{
                backgroundColor: "#374151",
                border: "none",
                borderRadius: "0.375rem",
              }} // ツールチップのスタイル
              itemStyle={{ color: "#E5E7EB" }} // テキストの色
              labelStyle={{ color: "#D1D5DB", fontWeight: "bold" }} // ラベルの色
              formatter={(value) => [
                `$${typeof value === "number" ? value.toFixed(2) : value}`,
                "価格",
              ]} // フォーマット調整
            />
            <defs>
              <linearGradient id="colorPrice" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#34D399" stopOpacity={0.8} />
                <stop offset="95%" stopColor="#34D399" stopOpacity={0.1} />
              </linearGradient>
            </defs>
            <Area
              type="monotone"
              dataKey="price"
              stroke="#34D399"
              fillOpacity={0.6}
              fill="url(#colorPrice)"
            />
          </AreaChart>
        </ResponsiveContainer>
      ) : (
        <div className="flex items-center justify-center h-full text-gray-500 dark:text-gray-400">
          チャートデータを表示できません。
        </div>
      )}
    </div>
  );
}

/**
 * 個別の株価情報を表示するカードコンポーネント
 * @param {object} props - コンポーネントのプロパティ
 * @param {object} props.stock - 株価データ
 * @param {boolean} props.isSelected - 選択中かどうか
 * @param {Function} props.onClick - カードクリック時の処理
 */
function StockCard({
  stock,
  isSelected,
  onClick,
}: {
  stock: {
    symbol: string;
    name: string;
    price: number;
    change: string;
    changePercent: string;
    volume: string;
    marketCap: string;
  };
  isSelected: boolean;
  onClick: (symbol: string) => void;
}) {
  const isPositive = stock.change.startsWith("+");
  const changeColor = isPositive ? "text-green-500" : "text-red-500";
  const borderClass = isSelected
    ? "border-blue-500 ring-2 ring-blue-500"
    : "border-gray-200 dark:border-gray-700";

  return (
    <div
      className={`bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md border ${borderClass} cursor-pointer transition-all duration-200 hover:shadow-lg hover:scale-[1.02]`}
      onClick={() => onClick(stock.symbol)}
      onKeyDown={() => onClick(stock.symbol)}
    >
      <div className="flex justify-between items-start mb-2">
        <div>
          <h4 className="text-lg font-bold text-gray-900 dark:text-white">{stock.symbol}</h4>
          <p className="text-xs text-gray-500 dark:text-gray-400 truncate w-32 sm:w-40">
            {stock.name}
          </p>
        </div>
        <div className="text-right">
          <p className="text-xl font-semibold text-gray-900 dark:text-white">${stock.price}</p>
          <p className={`text-sm font-medium ${changeColor}`}>
            {stock.change} ({stock.changePercent})
          </p>
        </div>
      </div>
      <div className="text-xs text-gray-500 dark:text-gray-400 flex justify-between mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
        <span>出来高: {stock.volume}</span>
        <span>時価総額: {stock.marketCap}</span>
      </div>
    </div>
  );
}

/**
 * ページネーションコントロール
 * @param {object} props
 * @param {number} props.currentPage - 現在のページ
 * @param {number} props.totalPages - 総ページ数
 * @param {Function} props.onPageChange - ページ変更時のハンドラ
 */
function Pagination({
  currentPage,
  totalPages,
  onPageChange,
}: {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  if (totalPages <= 1) return null; // 1ページ以下の場合は表示しない

  return (
    <div className="flex justify-center items-center space-x-2 mt-6">
      <button
        type="button"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors hover:bg-gray-300 dark:hover:bg-gray-600"
      >
        前へ
      </button>
      <span className="text-gray-700 dark:text-gray-300">
        {currentPage} / {totalPages}
      </span>
      <button
        type="button"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors hover:bg-gray-300 dark:hover:bg-gray-600"
      >
        次へ
      </button>
    </div>
  );
}

/**
 * ダッシュボード全体のページコンポーネント
 */
function DashboardPage() {
  const [allStocks] = useState(initialStockData); // 全データは変更しない
  const [selectedSymbol, setSelectedSymbol] = useState(
    initialStockData.length > 0 ? initialStockData[0].symbol : null,
  ); // 初期選択
  const [searchQuery, setSearchQuery] = useState(""); // 検索クエリ
  const [currentPage, setCurrentPage] = useState(1); // 現在のページ番号

  // 検索フィルタリング
  const filteredStocks = useMemo(() => {
    const lowerCaseQuery = searchQuery.toLowerCase();
    return allStocks.filter(
      (stock) =>
        stock.symbol.toLowerCase().includes(lowerCaseQuery) ||
        stock.name.toLowerCase().includes(lowerCaseQuery),
    );
  }, [allStocks, searchQuery]);

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
  const currentChartData = useMemo(() => {
    // ここで selectedSymbol に基づいて適切なチャートデータを返す
    // 今回はダミーデータをそのまま使う
    console.log("Updating chart for:", selectedSymbol);
    return chartData;
  }, [selectedSymbol]);

  return (
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
                isSelected={stock.symbol === selectedSymbol}
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
        <StockChart data={currentChartData} selectedSymbol={selectedSymbol} />
      </div>

      {/* フッター等、他の要素をここに追加可能 */}
      <footer className="mt-8 text-center text-gray-500 dark:text-gray-400 text-sm">
        データはダミーです。実際の取引には使用しないでください。
      </footer>
    </div>
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
