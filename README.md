# 株価通知システム (Stock Notification System)

リアルタイムの株価データを取得・保存し、メール通知機能とWebダッシュボードを提供する株式市場監視システムです。

## 🚀 主な機能

- **リアルタイム株価取得**: Yahoo Financeから主要株価指数のデータを自動取得
- **自動メール通知**: MailerSendを使用した日次市場サマリー配信
- **GraphQL API**: 株価データの効率的なクエリとAPI提供
- **Webダッシュボード**: インタラクティブなチャートと検索機能
- **OAuth認証**: セキュアなユーザーセッション管理
- **Docker対応**: 完全なコンテナ化環境
- **CLI ツール**: 運用・メンテナンス用コマンドラインインターフェース

## 📊 対応市場

- **日経平均株価** (^N225)
- **S&P 500** (^GSPC)
- **ダウ・ジョーンズ** (^DJI)
- **NASDAQ** (^IXIC)
- **カスタム指数** (^XDN)

## 🛠️ 技術スタック

### バックエンド (API)
- **Go 1.24+** - メインアプリケーション言語
- **GraphQL** - APIレイヤー (gqlgen)
- **PostgreSQL** - データベース
- **Docker** - コンテナ化
- **Wire** - 依存性注入
- **Cobra** - CLIフレームワーク

### フロントエンド (Web)
- **React 19** - UIフレームワーク
- **React Router v7** - SSRとルーティング
- **TypeScript** - 型安全性
- **TailwindCSS v4** - スタイリング
- **Apollo Client** - GraphQLクライアント
- **Recharts** - チャートライブラリ
- **Bun** - パッケージマネージャー

## 🚀 クイックスタート

### 前提条件

- Docker & Docker Compose
- Go 1.24+ (開発時)
- Bun または Node.js 20+ (フロントエンド開発時)

### 1. プロジェクトのクローン

```bash
git clone https://github.com/heyjun3/notify-stock.git
cd notify-stock
```

### 2. 環境変数の設定

```bash
# APIサーバー用環境変数
export DB_HOST="localhost"
export DB_PORT="5555"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="notify-stock"
export MAIL_TOKEN="your-mailersend-token"
export FROM="sender@example.com"
export TO="recipient@example.com"
export OAUTH_CLIENT_ID="your-oauth-client-id"
export OAUTH_CLIENT_SECRET="your-oauth-secret"
export OAUTH_REDIRECT_URL="http://localhost:8080/auth/callback"

# Webアプリ用環境変数
export VITE_BACKEND_URL="http://localhost:8080/"
```

### 3. データベースセットアップ

```bash
cd api
make db-setup
```

### 4. アプリケーションの起動

#### Docker Composeを使用 (推奨)

```bash
# APIサーバーとデータベースを起動
cd api
docker-compose up -d

# Webアプリケーションを起動
cd ../web
docker-compose up -d
```

#### 開発モードで起動

```bash
# APIサーバー
cd api
go run cmd/main.go server

# Webアプリケーション (別ターミナル)
cd web
bun install
bun run dev
```

### 5. アクセス

- **GraphQL Playground**: http://localhost:8080
- **Webダッシュボード**: http://localhost (Docker) または http://localhost:5173 (開発)

## 📋 使用方法

### 株価データの取得

```bash
cd api

# 過去7日分のデータを取得
go run cmd/main.go stock update

# 全履歴データを取得 (5年分)
go run cmd/main.go stock update -a

# 特定銘柄のみ更新
go run cmd/main.go stock update -s "^N225,^GSPC"
```

### メール通知の送信

```bash
# 日経平均とS&P500の通知を送信
go run cmd/main.go notify -s "^N225,^GSPC"

# 自動通知設定 (Makefile使用)
make notify
```

### GraphQL APIの使用

```graphql
# 全銘柄の取得
query {
  symbols {
    detail {
      symbol
      shortName
      price
      change
      changePercent
      volume
      marketCap
    }
  }
}

# 特定銘柄の価格履歴
query {
  symbol(symbol: "^GSPC") {
    detail {
      symbol
      shortName
      price
    }
    chart(input: {period: "1Y"}) {
      timestamp
      price
    }
  }
}

# 通知設定の作成 (認証が必要)
mutation {
  createNotification(input: {
    name: "My Portfolio"
    symbols: ["^GSPC", "^N225"]
  }) {
    id
    name
  }
}
```

## 🔧 開発

### API開発

```bash
cd api

# ビルド
make build

# GraphQLコード生成
make gqlgen

# テスト実行
go test ./internal/...

# データベース接続
make db-connect
```

### Web開発

```bash
cd web

# 開発サーバー起動
bun run dev

# ビルド
bun run build

# 型チェック
bun run typecheck

# コードフォーマット
bun run fmt && bun run lint

# GraphQL型生成
bun run codegen
```

### 新しい株式銘柄の追加

1. `api/config.yaml` に銘柄シンボルを追加
2. データベースの `symbols` テーブルにメタデータを挿入
3. `stock update` コマンドで履歴データを取得

### GraphQLスキーマの変更

```bash
cd api
# schema.graphqls を編集後
make gqlgen

cd ../web
# 型情報を更新
bun run codegen
```

## 🐳 本番デプロイ

### Docker Composeでの全体デプロイ

```bash
# 環境変数ファイルを作成
cp api/.env.example api/.env
cp web/.env.example web/.env

# 本番環境で起動
docker-compose -f api/compose.yaml -f web/compose.yaml up -d
```

### 個別コンテナビルド

```bash
# APIサーバー
cd api
docker build -t notify-stock-api .
docker run -p 8080:8080 notify-stock-api

# Webアプリ
cd web
docker build -f docker/Dockerfile -t notify-stock-web .
docker run -p 80:80 notify-stock-web
```

## 📊 アーキテクチャ

### システム構成

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │───▶│   GraphQL API   │───▶│   PostgreSQL    │
│   (React App)   │    │   (Go Server)   │    │   (Database)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       ▼                       │
         │            ┌─────────────────┐                │
         │            │  Yahoo Finance  │                │
         │            │   (Data Source) │                │
         │            └─────────────────┘                │
         │                       │                       │
         │                       ▼                       │
         │            ┌─────────────────┐                │
         └───────────▶│   MailerSend    │◀───────────────┘
                      │   (Email API)   │
                      └─────────────────┘
```

### ディレクトリ構造

```
notify-stock/
├── api/                    # Go バックエンド
│   ├── cmd/               # CLI コマンド
│   ├── internal/          # ビジネスロジック
│   ├── graph/             # GraphQL スキーマ・リゾルバー
│   ├── config.yaml        # 銘柄設定
│   └── compose.yaml       # Docker 設定
├── web/                   # React フロントエンド
│   ├── app/              # アプリケーションコード
│   ├── docker/           # Docker 設定
│   └── compose.yaml      # Docker 設定
└── README.md             # このファイル
```

## 🔐 環境変数

### API サーバー

| 変数名 | 説明 | 必須 | デフォルト |
|--------|------|------|-----------|
| `DB_HOST` | データベースホスト | No | localhost |
| `DB_PORT` | データベースポート | No | 5555 |
| `DB_USER` | データベースユーザー | No | postgres |
| `DB_PASSWORD` | データベースパスワード | No | postgres |
| `DB_NAME` | データベース名 | No | notify-stock |
| `MAIL_TOKEN` | MailerSend APIトークン | Yes | - |
| `FROM` | 送信者メールアドレス | Yes | - |
| `TO` | 受信者メールアドレス | Yes | - |
| `OAUTH_CLIENT_ID` | OAuth クライアントID | Yes | - |
| `OAUTH_CLIENT_SECRET` | OAuth クライアントシークレット | Yes | - |
| `OAUTH_REDIRECT_URL` | OAuth リダイレクトURL | Yes | - |

### Web アプリケーション

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `VITE_BACKEND_URL` | GraphQL APIのベースURL | http://localhost:8080/ |

## 📝 今後の予定

- [ ] 本番環境へのデプロイ
- [ ] ログイン機能の強化
- [ ] S&P500(JPY) 対応
- [ ] 移動平均線の表示機能
- [ ] テストカバレッジの向上
- [ ] データベースマイグレーション

## 🤝 コントリビューション

1. フォークしてください
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。

---

**注意**: このアプリケーションはデモンストレーション・学習目的です。実際の投資判断には使用しないでください。
