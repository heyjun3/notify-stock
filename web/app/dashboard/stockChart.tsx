import {
  ResponsiveContainer,
  AreaChart,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
  Area,
} from "recharts";

type StockChartProps = {
  data: { name: string; price: number }[]; // チャートデータ
  selectedSymbol: string | null; // 選択中の銘柄シンボル
};

/**
 * 株価チャートを表示するコンポーネント
 * @param {object} props - コンポーネントのプロパティ
 * @param {Array<object>} props.data - チャートデータ
 * @param {string} props.selectedSymbol - 選択中の銘柄シンボル
 */
export function StockChart({ data, selectedSymbol }: StockChartProps) {
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
              interval={"preserveEnd"} // ラベルを常に表示（データ数が多い場合は調整）
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
