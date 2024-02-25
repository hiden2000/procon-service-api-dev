package commonerrors

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

// AccessDeniedError - アクセス拒否エラー
type AccessDeniedError struct {
	Message string // エラーの詳細メッセージ
}

func (e *AccessDeniedError) Error() string {
	return e.Message
}

// NewAccessDeniedError - 新しいアクセス拒否エラーを生成
func NewAccessDeniedError(message string) *AccessDeniedError {
	return &AccessDeniedError{
		Message: message,
	}
}

// TokenError - JWTトークン関連のエラー
type TokenError struct {
	Type    string // エラーのタイプを示す（"InvalidToken", "TokenExpired" など）
	Message string // エラーの詳細メッセージ
}

func (e *TokenError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// 新しいTokenErrorを生成するヘルパー関数
func NewTokenError(errorType, message string) *TokenError {
	return &TokenError{
		Type:    errorType,
		Message: message,
	}
}

// WrapTokenError - JWTトークン操作中に発生したエラーをTokenErrorにラップする
func WrapTokenError(err error) *TokenError {
	if err == nil {
		return nil
	}

	// JWTライブラリ固有のエラータイプに基づいて条件分岐
	var ve *jwt.ValidationError
	if errors.As(err, &ve) { // JWTのバリデーションエラー
		switch {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return NewTokenError("InvalidTokenFormat", "Invalid token format")
		case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
			return NewTokenError("TokenExpired", "Token expired or not active yet")
		default:
			return NewTokenError("TokenValidationFailed", "Token validation failed")
		}
	}

	// RedisのNilエラーの場合
	if errors.Is(err, redis.Nil) {
		return NewTokenError("SessionNotFound", "Session not found or expired")
	}

	// その他のエラー
	return NewTokenError("UnexpectedError", "Unexpected Error occurred")
}
