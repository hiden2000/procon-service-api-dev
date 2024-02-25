package commonerrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

// APIError - APIレスポンス用のエラー
type APIError struct {
	StatusCode int    // HTTPステータスコード
	Message    string // クライアントに返すメッセージ
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError - 新しいAPIErrorを生成
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// ErrorToAPIError - エラーをAPIErrorに変換
func ErrorToAPIError(err error) *APIError {
	if err == nil {
		return nil
	}
	// エラータイプに応じた処理
	switch e := err.(type) {
	case *DBError:
		return NewAPIError(http.StatusInternalServerError, "Internal server error")
	case *NotFoundError:
		return NewAPIError(http.StatusNotFound, fmt.Sprintf("%s not found", e.Resource))
	case *ValidationError:
		return NewAPIError(http.StatusBadRequest, e.Error())
	case *ForeignKeyViolationError:
		return NewAPIError(http.StatusBadRequest, "Invalid reference")
	case *ConstraintViolationError:
		return NewAPIError(http.StatusBadRequest, e.Error())
	case *TransactionError:
		return NewAPIError(http.StatusInternalServerError, "Transaction error")
	case *TokenError:
		if e.Type == "InvalidTokenFormat" || e.Type == "TokenValidationFailed" {
			return NewAPIError(http.StatusUnauthorized, e.Message)
		} else if e.Type == "TokenExpired" {
			return NewAPIError(http.StatusUnauthorized, "Token expired")
		}
		return NewAPIError(http.StatusInternalServerError, e.Message)
	case *RequestParsingError, *RequestVariableError:
		return NewAPIError(http.StatusBadRequest, e.Error())
	case *MinIOError:
		return NewAPIError(http.StatusInternalServerError, "File storage error")
	case *FileValidationError:
		return NewAPIError(http.StatusBadRequest, e.Error())
	case *AccessDeniedError:
		return NewAPIError(http.StatusForbidden, e.Error())
	case *DataMismatchError:
		return NewAPIError(http.StatusBadRequest, "Invalid data format")
	default:
		// JWTのValidationErrorを特定する
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			// JWTのバリデーションエラーに基づいて適切なAPIErrorを生成
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return NewAPIError(http.StatusUnauthorized, "Invalid token format")
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return NewAPIError(http.StatusUnauthorized, "Token expired")
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return NewAPIError(http.StatusUnauthorized, "Token not active yet")
			default:
				return NewAPIError(http.StatusUnauthorized, "Token validation failed")
			}
		}
		// RedisのNilエラーを特定する
		if errors.Is(err, redis.Nil) {
			return NewAPIError(http.StatusUnauthorized, "Session not found or expired")
		}
		// その他のエラーは一般的なエラーメッセージで処理
		return NewAPIError(http.StatusInternalServerError, "Unexpected error occurred")
	}
}
