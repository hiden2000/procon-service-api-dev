package commonerrors

import "fmt"

// MinIOError - MinIO操作関連のエラー
type MinIOError struct {
	Action  string // 実行しようとしたアクション
	Message string // エラーの詳細メッセージ
}

func (e *MinIOError) Error() string {
	return fmt.Sprintf("MinIO error on %s: %s", e.Action, e.Message)
}

// NewMinIOError - MinIOErrorの生成
func NewMinIOError(action, message string) *MinIOError {
	return &MinIOError{
		Action:  action,
		Message: message,
	}
}

// WrapMinIOError - MinIO操作中に発生したエラーをMinIOErrorにラップする
func WrapMinIOError(action string, err error) *MinIOError {
	if err == nil {
		return nil
	}
	return NewMinIOError(action, err.Error())
}

// FileValidationError - ファイル検証エラー
type FileValidationError struct {
	Message string // エラーの詳細メッセージ
}

func (e *FileValidationError) Error() string {
	return fmt.Sprintf("File validation error: %s", e.Message)
}

// NewFileValidationError - FileValidationErrorの生成
func NewFileValidationError(message string) *FileValidationError {
	return &FileValidationError{
		Message: message,
	}
}
