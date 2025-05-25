# 株価ダッシュボード (Stock Price Dashboard)

株価情報を可視化するモダンなWebアプリケーション。リアルタイムの株価データ、インタラクティブなチャート、検索・フィルタリング機能を提供します。

## 🚀 特徴

- **リアルタイム株価表示**: 株価、変動率、出来高、時価総額を含む詳細な株価情報
- **インタラクティブチャート**: 1ヶ月、6ヶ月、1年、5年の期間選択可能な価格推移チャート  
- **検索・フィルタリング**: 銘柄シンボルや会社名による高速検索
- **ページネーション**: 大量の銘柄データを効率的に表示
- **レスポンシブデザイン**: モバイル・デスクトップ対応、ダークモード対応
- **サーバーサイドレンダリング**: React Router v7による高速な初期ロード

## 🛠️ 技術スタック

- **フロントエンド**: React 19, React Router v7, TypeScript
- **スタイリング**: TailwindCSS v4
- **データ取得**: Apollo Client, GraphQL
- **チャート**: Recharts
- **ビルドツール**: Vite, Bun
- **コード品質**: Biome (フォーマット・リント)
- **デプロイ**: Docker, Nginx

## 📋 前提条件

- Bun (推奨) または Node.js 20+
- GraphQL APIサーバー (バックエンド) が `http://localhost:8080/query` で実行中

## 🚀 クイックスタート

### 1. 依存関係のインストール

```bash
bun install
```

### 2. 環境変数の設定

```bash
# .env.local ファイルを作成
VITE_BACKEND_URL=http://localhost:8080/
```

### 3. GraphQL型の生成

```bash
bun run codegen
```

### 4. 開発サーバーの起動

```bash
bun run dev
```

アプリケーションは `http://localhost:5173` でアクセス可能です。

## 📝 開発コマンド

```bash
# 開発
bun run dev              # 開発サーバー起動 (HMR有効)

# ビルド・型チェック
bun run build           # プロダクションビルド
bun run typecheck       # TypeScript型チェック

# コード品質
bun run fmt             # Biomeによるフォーマット
bun run lint            # Biomeによるリント・自動修正

# GraphQL
bun run codegen         # GraphQLスキーマから型生成

# プロダクション
bun run start           # プロダクションサーバー起動
```

## 🐳 Docker デプロイ

### ローカルでのビルド・実行

```bash
# Dockerイメージのビルド
docker build -t stock-dashboard .

# コンテナの実行
docker run -p 80:80 stock-dashboard
```

### Docker Composeでの実行

```bash
# バックエンドAPIと連携して実行
docker-compose up
```

アプリケーションは `http://localhost` でアクセス可能です。

## 📁 プロジェクト構造

```
web/
├── app/                    # アプリケーションコード
│   ├── dashboard/          # ダッシュボード関連コンポーネント
│   │   ├── dashboard.tsx   # メインダッシュボード
│   │   ├── stockCard.tsx   # 株価カード
│   │   ├── stockChart.tsx  # 価格チャート
│   │   ├── pagination.tsx  # ページネーション
│   │   └── getSymbol.gql   # GraphQLクエリ
│   ├── gen/                # 自動生成ファイル
│   │   └── graphql.ts      # GraphQL型・フック
│   ├── routes/             # ルート定義
│   └── root.tsx            # ルートコンポーネント
├── docker/                 # Docker設定
├── public/                 # 静的ファイル
└── build/                  # ビルド出力
```

## 🔧 GraphQL スキーマ

アプリケーションは以下のGraphQLクエリを使用します:

```graphql
query GetSymbols($chartInput: ChartInput!) {
  symbols {
    detail {
      symbol          # 銘柄シンボル (例: AAPL)
      shortName       # 短縮名 (例: Apple Inc.)
      longName        # 正式名称
      price           # 現在価格
      change          # 価格変動
      changePercent   # 変動率
      volume          # 出来高
      marketCap       # 時価総額
      currencySymbol  # 通貨記号
    }
    chart(input: $chartInput) {
      timestamp       # タイムスタンプ
      price          # 価格
    }
  }
}
```

## 🌟 主要機能

### 株価ダッシュボード
- 株価カードによる一覧表示
- リアルタイム価格更新
- 変動率による色分け表示

### インタラクティブチャート
- 期間選択 (1M/6M/1Y/5Y)
- マウスホバーでの詳細表示
- レスポンシブデザイン

### 検索・フィルタリング
- 銘柄シンボル検索
- 会社名検索
- リアルタイムフィルタリング

## 🔒 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `VITE_BACKEND_URL` | GraphQL APIのベースURL | `http://localhost:8080/` |

## 📚 開発ガイド

### 新しいコンポーネントの追加
1. `app/dashboard/` ディレクトリに `.tsx` ファイルを作成
2. 既存のコンポーネントの命名規則とスタイルに従う
3. TypeScript型を適切に定義

### GraphQLクエリの変更
1. `.gql` ファイルを編集
2. `bun run codegen` を実行して型を再生成
3. 必要に応じてコンポーネントの型を更新

### スタイリング
- TailwindCSSユーティリティクラスを使用
- ダークモード対応 (`dark:` プレフィックス)
- レスポンシブデザイン (`sm:`, `lg:` プレフィックス)

## 🤝 コントリビューション

1. フォークしてください
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

---

**注意**: このアプリケーションはデモンストレーション目的です。実際の取引には使用しないでください。