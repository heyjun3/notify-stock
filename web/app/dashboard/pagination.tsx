type PaginationProps = {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
};

/**
 * ページネーションコントロール
 * @param {object} props
 * @param {number} props.currentPage - 現在のページ
 * @param {number} props.totalPages - 総ページ数
 * @param {Function} props.onPageChange - ページ変更時のハンドラ
 */
export function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
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
