export type Period = "1M" | "6M" | "1Y" | "5Y"; // 期間の型定義
const periods: Period[] = ["1M", "6M", "1Y", "5Y"] as const; // 期間のリスト

type PeriodSelectorProps = {
  currentPeriod: Period;
  onPeriodChange: (period: Period) => void;
};

/**
 * 期間選択ボタン
 * @param {object} props
 * @param {string} props.currentPeriod - 現在選択中の期間キー ('1M', '6M', etc.)
 * @param {Function} props.onPeriodChange - 期間変更時のハンドラ
 */
export function PeriodSelector({ currentPeriod, onPeriodChange }: PeriodSelectorProps) {
  return (
    <div className="flex justify-end space-x-1 sm:space-x-2 mb-4">
      {periods.map((period) => (
        <button
          key={period}
          onClick={() => onPeriodChange(period)}
          className={`px-3 py-1 rounded-md text-xs sm:text-sm transition-colors ${
            currentPeriod === period
              ? "bg-blue-500 text-white font-semibold"
              : "bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600"
          }`}
        >
          {period} {/* ラベルを短く表示 */}
        </button>
      ))}
    </div>
  );
}
