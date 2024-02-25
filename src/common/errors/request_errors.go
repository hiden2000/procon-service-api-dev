package commonerrors

import "fmt"

// RequestParsingError - リクエスト解析エラー
type RequestParsingError struct {
	Message string // エラーメッセージ
}

func (e *RequestParsingError) Error() string {
	return e.Message
}

func NewRequestParsingError(msg string) *RequestParsingError {
	return &RequestParsingError{
		Message: msg,
	}
}

// RequestVariableError - リクエスト変数エラー
type RequestVariableError struct {
	VariableName string // 問題の変数名
	Message      string // エラーメッセージ
}

func (e *RequestVariableError) Error() string {
	return fmt.Sprintf("variable %s error: %s", e.VariableName, e.Message)
}

// WrapRequestParsingError - リクエストボディ解析エラーのラッピング
func WrapRequestParsingError(err error) *RequestParsingError {
	return &RequestParsingError{
		Message: "Failed to decode request body: " + err.Error(),
	}
}

// WrapRequestVariableError - リクエストからの変数取得エラーのラッピング
func WrapRequestVariableError(variableName string, err error) *RequestVariableError {
	return &RequestVariableError{
		VariableName: variableName,
		Message:      err.Error(),
	}
}
