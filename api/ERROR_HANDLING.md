# エラーハンドリング統一ガイド

## 概要

このプロジェクトでは統一されたエラーハンドリングパターンを採用しています。

## エラーの種類とコード

```go
const (
    ErrCodeInternalServer    ErrorCode = "INTERNAL_SERVER_ERROR"
    ErrCodeBadRequest        ErrorCode = "BAD_REQUEST"
    ErrCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
    ErrCodeNotFound          ErrorCode = "NOT_FOUND"
    ErrCodeValidation        ErrorCode = "VALIDATION_ERROR"
    ErrCodeExternalService   ErrorCode = "EXTERNAL_SERVICE_ERROR"
    ErrCodeDatabase          ErrorCode = "DATABASE_ERROR"
    ErrCodeSession           ErrorCode = "SESSION_ERROR"
    ErrCodeConfiguration     ErrorCode = "CONFIGURATION_ERROR"
)
```

## 使用例

### 1. バリデーションエラー

```go
func validateUser(user *User) error {
    if user.Email == "" {
        return NewValidationError("Email is required", "email field cannot be empty")
    }
    return nil
}
```

### 2. データベースエラーのラップ

```go
func getUserByID(db *bun.DB, id string) (*User, error) {
    user := &User{}
    err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, NewNotFoundError("User")
        }
        return nil, WrapError(err, ErrCodeDatabase, "Failed to fetch user")
    }
    return user, nil
}
```

### 3. 外部サービスエラー

```go
func fetchStockData(symbol string) (*StockData, error) {
    resp, err := http.Get(fmt.Sprintf("https://api.example.com/stock/%s", symbol))
    if err != nil {
        return nil, WrapError(err, ErrCodeExternalService, "Failed to fetch stock data")
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return nil, NewAppError(ErrCodeExternalService, 
            fmt.Sprintf("Stock API returned status %d", resp.StatusCode), nil)
    }
    
    // データ処理...
}
```

### 4. HTTPハンドラーでのエラーレスポンス

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        WriteErrorResponse(w, NewValidationError("User ID is required", "id parameter is missing"))
        return
    }
    
    user, err := getUserByID(db, userID)
    if err != nil {
        WriteErrorResponse(w, err)  // AppErrorはそのまま渡せる
        return
    }
    
    // 成功レスポンス...
}
```

## エラーレスポンス例

### 開発環境

```json
{
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": "email field cannot be empty"
}
```

### 本番環境（内部エラーの場合）

```json
{
    "code": "INTERNAL_SERVER_ERROR", 
    "message": "Internal server error occurred"
}
```

## ベストプラクティス

1. **エラーのラップ**: `WrapError()` を使って元のエラーを保持しつつコンテキストを追加
2. **適切なエラーコード**: 操作の種類に応じて適切なエラーコードを選択
3. **ログ出力**: `WriteErrorResponse()` が自動的にログ出力を行うため、手動でのログ出力は不要
4. **セキュリティ**: 本番環境では内部エラーの詳細を隠蔽
5. **一貫性**: 全てのHTTPエラーレスポンスで `WriteErrorResponse()` を使用

## 移行ガイド

### Before (旧方式)
```go
if err != nil {
    logger.Error("Database error", "error", err)
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
}
```

### After (新方式)
```go
if err != nil {
    WriteErrorResponse(w, WrapError(err, ErrCodeDatabase, "Failed to fetch data"))
    return
}
```