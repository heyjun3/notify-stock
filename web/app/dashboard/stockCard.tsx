type Stock = {
  symbol: string;
  shortName: string;
  longName: string;
  price: number;
  change: string;
  changePercent: string;
  volume?: string | null;
  marketCap?: string | null;
  currencySymbol: string;
};

type StockCardProps = {
  stock: Stock;
  isSelected: boolean;
  onClick: (symbol: { symbol: string; shortName: string }) => void;
};

/**
 * 個別の株価情報を表示するカードコンポーネント
 * @param {object} props - コンポーネントのプロパティ
 * @param {object} props.stock - 株価データ
 * @param {boolean} props.isSelected - 選択中かどうか
 * @param {Function} props.onClick - カードクリック時の処理
 */
export function StockCard({ stock, isSelected, onClick }: StockCardProps) {
  const isPositive = stock.change.startsWith("+");
  const changeColor = isPositive ? "text-green-500" : "text-red-500";
  const borderClass = isSelected
    ? "border-blue-500 ring-2 ring-blue-500"
    : "border-gray-200 dark:border-gray-700";

  return (
    <div
      className={`bg-white dark:bg-gray-800 p-4 rounded-lg shadow-md border ${borderClass} cursor-pointer transition-all duration-200 hover:shadow-lg hover:scale-[1.02]`}
      onClick={() => onClick(stock)}
      onKeyDown={() => onClick(stock)}
    >
      <div className="flex justify-between items-start mb-2">
        <div>
          <h4 className="text-lg font-bold text-gray-900 dark:text-white">{stock.shortName}</h4>
          <p className="text-xs text-gray-500 dark:text-gray-400 truncate w-32 sm:w-40">
            {stock.longName}
          </p>
        </div>
        <div className="text-right">
          <p className="text-xl font-semibold text-gray-900 dark:text-white">
            {stock.currencySymbol}
            {stock.price}
          </p>
          <p className={`text-sm font-medium ${changeColor}`}>
            {stock.change} ({stock.changePercent})
          </p>
        </div>
      </div>
      {stock.volume || stock.marketCap ? (
        <div className="text-xs text-gray-500 dark:text-gray-400 flex justify-between mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
          {stock.volume && <span>出来高: {stock.volume}</span>}
          {stock.marketCap && <span>時価総額: {stock.marketCap}</span>}
        </div>
      ) : (
        <div />
      )}
    </div>
  );
}
