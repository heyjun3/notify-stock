package notifystock

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// エラーコード定義
type ErrorCode string

const (
	ErrCodeInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeDatabase        ErrorCode = "DATABASE_ERROR"
	ErrCodeSession         ErrorCode = "SESSION_ERROR"
	ErrCodeConfiguration   ErrorCode = "CONFIGURATION_ERROR"
)

// アプリケーションエラー構造体
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// エラー生成ヘルパー関数
func NewAppError(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func NewValidationError(message string, details string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Details: details,
	}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Code:    ErrCodeInternalServer,
		Message: message,
		Cause:   cause,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
	}
}

// HTTPステータスコードマッピング
func (e *AppError) HTTPStatusCode() int {
	switch e.Code {
	case ErrCodeBadRequest, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeInternalServer, ErrCodeDatabase, ErrCodeExternalService:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// HTTPエラーレスポンス送信
func WriteErrorResponse(w http.ResponseWriter, err error) {
	var appErr *AppError

	// AppErrorにキャスト、失敗したら内部サーバーエラーとして扱う
	if !errors.As(err, &appErr) {
		appErr = NewInternalError("An unexpected error occurred", err)
	}

	// ログ出力（本番環境では詳細な情報をログに記録）
	if Cfg.IsProduction() {
		logger.Error("Error occurred",
			"code", appErr.Code,
			"message", appErr.Message,
			"cause", appErr.Cause)
	} else {
		logger.Error("Error occurred",
			"error", appErr.Error())
	}

	// HTTPレスポンス送信
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatusCode())

	// 本番環境では内部エラーの詳細を隠す
	response := appErr
	if Cfg.IsProduction() && appErr.Code == ErrCodeInternalServer {
		response = &AppError{
			Code:    ErrCodeInternalServer,
			Message: "Internal server error occurred",
		}
	}

	json.NewEncoder(w).Encode(response)
}

// エラーチェインの作成
func WrapError(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}
