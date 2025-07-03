import React, { useState } from "react";
import { Link } from "react-router";
import { Trash2, BellPlus, Clock, Briefcase, CheckSquare, Square, LogOut } from "lucide-react";
import { useCreateNotification } from "./hooks/createNotification";
import { useDeleteNotification } from "./hooks/deleteNotification";
import { useGetNotification } from "./hooks/getNotification";

type Stock = {
  symbol: string;
  name: string;
};

type NotificationFormProps = {
  stocks: Stock[];
  onAddNotification: (notification: {
    id: number;
    time: string;
    selectedStockSymbols: string[];
  }) => void;
};

function NotificationForm({ stocks, onAddNotification }: NotificationFormProps) {
  const [time, setTime] = useState("09:00");
  const [selectedStockSymbols, setSelectedStockSymbols] = useState<string[]>([]);
  const [error, setError] = useState("");

  const handleStockSelectionChange = (symbol: string) => {
    setSelectedStockSymbols((prevSelected) =>
      prevSelected.includes(symbol)
        ? prevSelected.filter((s) => s !== symbol)
        : [...prevSelected, symbol],
    );
  };

  const handleSubmit = (e: React.ChangeEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!time || selectedStockSymbols.length === 0) {
      setError("通知時間と対象の株（1つ以上）を選択してください。");
      return;
    }
    setError("");
    onAddNotification({ id: Date.now(), time, selectedStockSymbols });
  };

  return (
    <div className="bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-md">
      <h4 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white flex items-center">
        <BellPlus size={20} className="mr-2 text-blue-500" />
        株価通知を登録
      </h4>
      {error && <p className="text-red-500 text-sm mb-3">{error}</p>}
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label
            htmlFor="notify-time"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1 flex items-center"
          >
            <Clock size={16} className="mr-1 text-gray-500" /> 通知時間
          </label>
          <input
            type="time"
            id="notify-time"
            value={time}
            onChange={(e) => setTime(e.target.value)}
            className="w-full md:w-1/2 lg:w-1/3 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 flex items-center">
            <Briefcase size={16} className="mr-1 text-gray-500" /> 対象の株 (複数選択可)
          </label>
          <div className="max-h-60 overflow-y-auto space-y-2 p-3 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-50 dark:bg-gray-700">
            {stocks.map((stock) => (
              <label
                key={stock.symbol}
                className="flex items-center space-x-2 p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-600 cursor-pointer"
              >
                <input
                  type="checkbox"
                  checked={selectedStockSymbols.includes(stock.symbol)}
                  onChange={() => handleStockSelectionChange(stock.symbol)}
                  className="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:bg-gray-700 dark:border-gray-600"
                />
                <span className="text-sm text-gray-800 dark:text-gray-200">
                  {stock.symbol} - {stock.name}
                </span>
                {selectedStockSymbols.includes(stock.symbol) ? (
                  <CheckSquare size={16} className="text-blue-500" />
                ) : (
                  <Square size={16} className="text-gray-400 dark:text-gray-500" />
                )}
              </label>
            ))}
          </div>
        </div>
        <button
          type="submit"
          className="w-full sm:w-auto px-6 py-2 bg-blue-500 text-white font-semibold rounded-md shadow-sm hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 transition-colors flex items-center justify-center"
        >
          <BellPlus size={18} className="mr-2" />
          通知を登録する
        </button>
      </form>
    </div>
  );
}

type Notification = {
  id: string;
  time: string;
  tickers: string[];
};

type NotificationListProps = {
  notifications: Notification[];
  onDeleteNotification: (id: string) => void;
};

function NotificationList({ notifications, onDeleteNotification }: NotificationListProps) {
  if (notifications.length === 0) {
    return (
      <div className="mt-8 bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-md text-center">
        <p className="text-gray-500 dark:text-gray-400">登録済みの通知はありません。</p>
      </div>
    );
  }
  return (
    <div className="mt-8 bg-white dark:bg-gray-800 p-4 sm:p-6 rounded-lg shadow-md">
      <h4 className="text-lg font-semibold mb-4 text-gray-900 dark:text-white">登録済み通知一覧</h4>
      <ul className="space-y-3">
        {notifications.map((notification) => (
          <li
            key={notification.id}
            className="flex flex-col sm:flex-row justify-between items-start sm:items-center p-3 bg-gray-50 dark:bg-gray-700 rounded-md shadow-sm"
          >
            <div className="flex-grow mb-2 sm:mb-0">
              <p className="text-sm text-gray-800 dark:text-gray-200">
                <Clock size={14} className="inline mr-1 text-gray-500 dark:text-gray-400" /> 時間:{" "}
                <span className="font-medium">{notification.time}</span>
              </p>
              <p className="text-sm text-gray-800 dark:text-gray-200">
                <Briefcase size={14} className="inline mr-1 text-gray-500 dark:text-gray-400" />{" "}
                銘柄: <span className="font-medium">{notification.tickers.join(", ")}</span>
              </p>
            </div>
            <button
              onClick={() => onDeleteNotification(notification.id)}
              className="px-3 py-1 bg-red-500 text-white text-xs font-semibold rounded-md hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 transition-colors flex items-center"
            >
              <Trash2 size={14} className="mr-1" />
              削除
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}

type NotificationSectionProps = {
  user: {
    displayName: string;
    photoURL?: string;
  } | null;
  allStocks: Stock[];
};

export function NotificationSection({ user, allStocks }: NotificationSectionProps) {
  const { loading, notifications, unAuthorization, refetch } = useGetNotification();
  const { handleCreateNotification } = useCreateNotification(refetch);
  const { handleDeleteNotification } = useDeleteNotification(refetch);
  const isAuthorized = !unAuthorization;

  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md min-h-[250px] flex flex-col justify-center">
      {isAuthorized && user ? (
        <div>
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 border-b dark:border-gray-700 pb-4">
            <div className="flex items-center mb-3 sm:mb-0">
              {user.photoURL && (
                <img
                  src={user.photoURL}
                  alt="User Avatar"
                  className="w-9 h-9 rounded-full mr-3 border-2 border-blue-500"
                />
              )}
              <span className="text-sm text-gray-700 dark:text-gray-300">
                ようこそ, <span className="font-semibold">{user.displayName}</span> さん
              </span>
            </div>
            <Link
              to={new URL("/logout", import.meta.env.VITE_BACKEND_URL).toString()}
              className="w-full sm:w-auto px-4 py-2 bg-red-500 text-white text-sm font-semibold rounded-md hover:bg-red-600 transition-colors flex items-center justify-center"
            >
              <LogOut size={16} className="mr-2" />
              ログアウト
            </Link>
          </div>
          {!loading && (
            <div className="space-y-4">
              {notifications.length ? (
                <NotificationList
                  notifications={notifications}
                  onDeleteNotification={handleDeleteNotification}
                />
              ) : (
                <NotificationForm stocks={allStocks} onAddNotification={handleCreateNotification} />
              )}
            </div>
          )}
        </div>
      ) : (
        <div className="text-center">
          <h4 className="text-lg font-semibold mb-3 text-gray-900 dark:text-white">
            通知機能を利用するにはログイン
          </h4>
          <p className="text-gray-500 dark:text-gray-400 mb-6 text-sm max-w-sm mx-auto">
            Googleアカウントでログインして、気になる銘柄の株価通知を設定しましょう。
          </p>
          <Link
            to={new URL("/login", import.meta.env.VITE_BACKEND_URL).toString()}
            className="px-6 py-3 bg-white text-gray-700 font-semibold rounded-lg shadow-md hover:bg-gray-100 transition-all transform hover:scale-105 border border-gray-300 flex items-center justify-center mx-auto w-80"
          >
            <svg className="w-5 h-5 mr-3" viewBox="0 0 48 48" aria-hidden="true">
              <path
                fill="#FFC107"
                d="M43.611,20.083H42V20H24v8h11.303c-1.649,4.657-6.08,8-11.303,8c-6.627,0-12-5.373-12-12s5.373-12,12-12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C12.955,4,4,12.955,4,24s8.955,20,20,20s20-8.955,20-20C44,22.659,43.862,21.35,43.611,20.083z"
              ></path>
              <path
                fill="#FF3D00"
                d="M6.306,14.691l6.571,4.819C14.655,15.108,18.961,12,24,12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C16.318,4,9.656,8.337,6.306,14.691z"
              ></path>
              <path
                fill="#4CAF50"
                d="M24,44c5.166,0,9.86-1.977,13.409-5.192l-6.19-5.238C29.211,35.091,26.715,36,24,36c-5.222,0-9.519-3.108-11.182-7.482l-6.571,4.819C9.656,39.663,16.318,44,24,44z"
              ></path>
              <path
                fill="#1976D2"
                d="M43.611,20.083H42V20H24v8h11.303c-0.792,2.237-2.231,4.166-4.087,5.571l6.19,5.238C42.022,35.244,44,30.036,44,24C44,22.659,43.862,21.35,43.611,20.083z"
              ></path>
            </svg>
            Googleでログイン
          </Link>
        </div>
      )}
    </div>
  );
}
