import React, { useState } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, AreaChart, Area } from 'recharts';

// --- Mock Data ---
// ダミーの株価データ（実際のアプリケーションではAPIから取得）
const initialStockData = [
  { symbol: 'AAPL', name: 'Apple Inc.', price: 175.32, change: '+1.25', changePercent: '+0.72%', volume: '98.7M', marketCap: '2.75T' },
  { symbol: 'GOOGL', name: 'Alphabet Inc.', price: 2850.88, change: '-15.40', changePercent: '-0.54%', volume: '1.5M', marketCap: '1.90T' },
  { symbol: 'MSFT', name: 'Microsoft Corp.', price: 340.15, change: '+2.80', changePercent: '+0.83%', volume: '35.2M', marketCap: '2.53T' },
  { symbol: 'AMZN', name: 'Amazon.com, Inc.', price: 140.50, change: '-0.95', changePercent: '-0.67%', volume: '55.6M', marketCap: '1.43T' },
];

// ダミーのチャートデータ（選択された銘柄に応じて変化させる想定）
const chartData = [
  { name: '1月', price: 160 },
  { name: '2月', price: 165 },
  { name: '3月', price: 170 },
  { name: '4月', price: 168 },
  { name: '5月', price: 175 },
  { name: '6月', price: 180 },
  { name: '7月', price: 178 },
  { name: '8月', price: 185 },
  { name: '9月', price: 182 },
  { name: '10月', price: 175 },
  { name: '11月', price: 178 },
  { name: '12月', price: 175 },
];

// --- Components ---

/**
 * 株価チャートを表示するコンポーネント
 * @param {object} props - コンポーネントのプロパティ
 * @param {Array<object>} props.data - チャートデータ
 * @param {string} props.selectedSymbol - 選択中の銘柄シンボル
 */
function StockChart({ data, selectedSymbol }: {
    data: Array<{ name: string; price: number }>;
    selectedSymbol: string;
}) {
  return (
    <div className="bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-md h-[300px] sm:h-[400px]">
      <h3 className="text-lg sm:text-xl font-semibold mb-4 text-gray-800 dark:text-gray-200">{selectedSymbol} - 過去12ヶ月の価格推移</h3>
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart
          data={data}
          margin={{
            top: 5,
            right: 20, // 右側のスペースを確保
            left: 0,  // 左側のスペースを調整
            bottom: 20, // 下部のスペースを確保
          }}
        >
          <CartesianGrid strokeDasharray="3 3" stroke="#4B5563" /> {/* グリッド線の色を調整 */}
          <XAxis
             dataKey="name"
             tick={{ fill: '#9CA3AF' }} // 目盛りの色
             tickLine={{ stroke: '#9CA3AF' }} // 目盛り線の色
             axisLine={{ stroke: '#9CA3AF' }} // 軸線の色
             padding={{ left: 10, right: 10 }} // X軸の両端にパディングを追加
             angle={-30} // ラベルを斜めにする
             textAnchor="end" // ラベルのアンカーを調整
             height={50} // X軸の高さを確保
           />
          <YAxis
             tick={{ fill: '#9CA3AF' }}
             tickLine={{ stroke: '#9CA3AF' }}
             axisLine={{ stroke: '#9CA3AF' }}
             domain={['auto', 'auto']} // Y軸の範囲を自動調整
             tickFormatter={(value) => `$${value}`} // ドル記号を追加
           />
          <Tooltip
            contentStyle={{ backgroundColor: '#374151', border: 'none', borderRadius: '0.375rem' }} // ツールチップのスタイル
            itemStyle={{ color: '#E5E7EB' }} // テキストの色
            labelStyle={{ color: '#D1D5DB', fontWeight: 'bold' }} // ラベルの色
            formatter={(value) => [`$${typeof value === 'number' ? value.toFixed(2) : value}`, "価格"]} // フォーマット調整
          />
          <Area type="monotone" dataKey="price" stroke="#34D399" fillOpacity={0.6} fill="url(#colorPrice)" />
           <defs>
             <linearGradient id="colorPrice" x1="0" y1="0" x2="0" y2="1">
               <stop offset="5%" stopColor="#34D399" stopOpacity={0.8}/>
               <stop offset="95%" stopColor="#34D399" stopOpacity={0.1}/>
             </linearGradient>
           </defs>
        </AreaChart>
      </ResponsiveContainer>
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
function StockCard({ stock, isSelected, onClick }: {
    stock: { symbol: string; name: string; price: number; change: string; changePercent: string; volume: string; marketCap: string };
    isSelected: boolean;
    onClick: (symbol: string) => void;
}) {
  const isPositive = stock.change.startsWith('+');
  const changeColor = isPositive ? 'text-green-500' : 'text-red-500';
  const borderClass = isSelected ? 'border-blue-500 ring-2 ring-blue-500' : 'border-gray-200 dark:border-gray-700';

  return (
    <div
      className={`bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md border ${borderClass} cursor-pointer transition-all duration-200 hover:shadow-lg hover:scale-[1.02]`}
      onClick={() => onClick(stock.symbol)}
    >
      <div className="flex justify-between items-start mb-2">
        <div>
          <h4 className="text-lg font-bold text-gray-900 dark:text-white">{stock.symbol}</h4>
          <p className="text-xs text-gray-500 dark:text-gray-400 truncate w-32 sm:w-40">{stock.name}</p>
        </div>
        <div className="text-right">
          <p className="text-xl font-semibold text-gray-900 dark:text-white">${stock.price}</p>
          <p className={`text-sm font-medium ${changeColor}`}>{stock.change} ({stock.changePercent})</p>
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
 * ダッシュボード全体のページコンポーネント
 */
function DashboardPage() {
  const [stocks, setStocks] = useState(initialStockData);
  const [selectedSymbol, setSelectedSymbol] = useState(initialStockData[0].symbol); // 初期選択

  // カードクリック時のハンドラ
  const handleCardClick = (symbol: string) => {
    setSelectedSymbol(symbol);
    // ここで選択された銘柄に応じてチャートデータを更新するロジックを追加
    // 例: fetchChartData(symbol).then(data => setChartData(data));
    console.log(`Selected stock: ${symbol}`);
  };

  return (
    <div className="p-4 sm:p-8 bg-gray-100 dark:bg-gray-900 min-h-screen font-sans">
      <h1 className="text-2xl sm:text-3xl font-bold mb-6 text-gray-900 dark:text-white">株価ダッシュボード</h1>

      {/* 株価カード一覧 */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {stocks.map((stock) => (
          <StockCard
            key={stock.symbol}
            stock={stock}
            isSelected={stock.symbol === selectedSymbol}
            onClick={handleCardClick}
          />
        ))}
      </div>

      {/* 株価チャート */}
      <StockChart data={chartData} selectedSymbol={selectedSymbol} />

      {/* フッター等、他の要素をここに追加可能 */}
      <footer className="mt-8 text-center text-gray-500 dark:text-gray-400 text-sm">
        データはダミーです。実際の取引には使用しないでください。
      </footer>
    </div>
  );
}

// --- Main App Component ---
// React Router を使う場合は、ここで <BrowserRouter>, <Routes>, <Route> を設定します。
// 今回はシンプルにするため、DashboardPage を直接表示します。
export default function App() {
  // ダークモード切り替え等の全体的な状態管理をここで行うことも可能
  // const [isDarkMode, setIsDarkMode] = useState(false);
  // const toggleDarkMode = () => setIsDarkMode(!isDarkMode);

  // className={isDarkMode ? 'dark' : ''} をルート要素に追加

  return (
    // <div className={isDarkMode ? 'dark' : ''}>
    <div>
      {/* ヘッダーやナビゲーションバーをここに追加可能 */}
      {/* <button onClick={toggleDarkMode}>Toggle Dark Mode</button> */}
      <DashboardPage />
    </div>
  );
}
