package commonerrors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// DBError - データベース操作関連のエラー
type DBError struct {
	Op  string // 操作（"insert", "update", "delete" 等）
	Err error  // 元となったエラー
}

func (e *DBError) Error() string {
	return fmt.Sprintf("db operation %s failed: %v", e.Op, e.Err)
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// NewDBError - DBErrorを生成するヘルパー関数
func NewDBError(op string, errMsg string) *DBError {
	return &DBError{
		Op:  op,
		Err: errors.New(errMsg),
	}
}

// NotFoundError - レコードが見つからない場合のエラー
type NotFoundError struct {
	Resource        string // リソースの種類（"User", "Problem", "Solution" 等）
	Identifier      string // 識別子の名前（"ID", "Username", "Email" 等）
	IdentifierValue string // 識別子の値
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with %s %s not found", e.Resource, e.Identifier, e.IdentifierValue)
}

// NewNotFoundError - NotFoundErrorを生成するヘルパー関数
func NewNotFoundError(resource, identifier, value string) *NotFoundError {
	return &NotFoundError{
		Resource:        resource,
		Identifier:      identifier,
		IdentifierValue: value,
	}
}

// TransactionError - トランザクション処理中に発生したエラー
type TransactionError struct {
	Stage string // エラーが発生した処理の段階（"begin", "commit", "rollback" など）
	Err   error  // 元となるエラー
}

func (e *TransactionError) Error() string {
	return "transaction " + e.Stage + " error: " + e.Err.Error()
}

func (e *TransactionError) Unwrap() error {
	return e.Err
}

// NewTransactionError - TransactionErrorを生成するヘルパー関数
func NewTransactionError(stage string, err error) *TransactionError {
	return &TransactionError{
		Stage: stage,
		Err:   err,
	}
}

// ConstraintViolationError - 制約違反に関するエラー
type ConstraintViolationError struct {
	Constraint string // 違反した制約の名前
	Err        error  // 元となるエラー
}

func (e *ConstraintViolationError) Error() string {
	return fmt.Sprintf("constraint violation: %s, %v", e.Constraint, e.Err)
}

// NewConstraintViolationError - ConstraintViolationErrorを生成するヘルパー関数
func NewConstraintViolationError(constraint, errMsg string) *ConstraintViolationError {
	return &ConstraintViolationError{
		Constraint: constraint,
		Err:        errors.New(errMsg),
	}
}

// ValidationError - データ検証エラー
type ValidationError struct {
	Field string // エラーのあったフィールド
	Err   error  // 元となるエラー
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s, %v", e.Field, e.Err)
}

func NewValidationError(field, errMsg string) *ValidationError {
	return &ValidationError{
		Field: field,
		Err:   errors.New(errMsg),
	}
}

// ForeignKeyViolationError - 外部キー制約違反エラー
type ForeignKeyViolationError struct {
	ForeignKey string // 違反した外部キーの名前
	Err        error  // 元となるエラー
}

func (e *ForeignKeyViolationError) Error() string {
	return fmt.Sprintf("foreign key violation: %s, %v", e.ForeignKey, e.Err)
}

// NewForeignKeyViolationError - ForeignKeyViolationErrorを生成するヘルパー関数
func NewForeignKeyViolationError(foreignKey, errMsg string) *ForeignKeyViolationError {
	return &ForeignKeyViolationError{
		ForeignKey: foreignKey,
		Err:        errors.New(errMsg),
	}
}

// DataMismatchError - データ型不一致エラー
type DataMismatchError struct {
	Err error // エラーの詳細メッセージ
}

func (e *DataMismatchError) Error() string {
	return e.Err.Error()
}

// NewDataMismatchError - DataMismatchErrorを生成するヘルパー関数
func NewDataMismatchError(errMsg string) *DataMismatchError {
	return &DataMismatchError{
		Err: errors.New(errMsg),
	}
}

func WrapDBError(op string, err error) error {
	if err == nil {
		return nil
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062: // ER_DUP_ENTRY
			return handleDuplicateEntryError(mysqlErr)
		case 1406, 3819: // CHECK制約違反
			return handleCheckConstraintError(mysqlErr)
		case 1452: // ER_NO_REFERENCED_ROW_2
			return handleForeignKeyError(mysqlErr)
		case 1366: // ER_TRUNCATED_WRONG_VALUE_FOR_FIELD
			return handleTruncatedValueError(mysqlErr)
		// TODO: 他のエラーケースを追加
		default:
			return handleUnexpectedDBError(op, err)
		}
	}
	// その他のエラー
	return NewDBError(op, err.Error())
}

// handleTruncatedValueError - データ形式不正によるエラーのハンドリング
func handleTruncatedValueError(err *mysql.MySQLError) error {
	// err.Messageから詳細を解析し，適切なカスタムエラーを返す
	// NOTE: 具体的なフィールドを特定するのが難しいので，一般的なエラーメッセージを使用
	return NewDataMismatchError("The format of the input data is invalid.")
}

// handleDuplicateEntryError - ER_DUP_ENTRYエラーのハンドリング
func handleDuplicateEntryError(err *mysql.MySQLError) error {
	// err.Messageから詳細を解析し，適切なカスタムエラーを返す
	// キー名に基づいて条件分岐
	// ユニーク制約違反の場合，どのフィールドが原因かを判断
	if strings.Contains(err.Message, "username_unique") {
		return NewConstraintViolationError("username", "This username is already in use.")
	}
	// 他のユニーク制約違反をチェック
	return NewDBError("INSERT", "A unique constraint violation occurred.")
}

// handleCheckConstraintError - CHECK制約違反のハンドリング
func handleCheckConstraintError(err *mysql.MySQLError) error {
	// err.Messageから詳細を解析し，適切なカスタムエラーを返す
	if strings.Contains(err.Message, "chk_Problems_Difficulty") {
		return NewValidationError("Difficulty", "The difficulty value is out of range (integer value between 1 - 5).")
	}
	if strings.Contains(err.Message, "Data too long for column 'Username'") {
		return NewValidationError("Username", "The Username length must be less than 255 characters.")
	}
	return NewDBError("CONSTRAINT", "A data constraint violation occurred.")
}

// handleForeignKeyError - 外部キー制約違反のハンドリング
func handleForeignKeyError(err *mysql.MySQLError) error {
	// err.Messageから詳細を解析し，適切なカスタムエラーを返す
	if strings.Contains(err.Message, "fk_Solutions_Users") {
		return NewForeignKeyViolationError("UserID", "The specified user does not exist.")
	}
	if strings.Contains(err.Message, "fk_Solutions_Problems") {
		return NewForeignKeyViolationError("ProblemID", "The specified problem does not exist.")
	}
	// 他の外部キー制約違反をチェック
	return NewDBError("FOREIGN_KEY", "A foreign key constraint violation occurred.")
}

// handleUnexpectedDBError - 予期しないデータベースエラーのハンドリング
func handleUnexpectedDBError(op string, err error) error {
	// 予期しないエラーであることを明示しつつ，詳細情報の露出を避ける
	return NewDBError(op, "An unexpected error occurred during database operations.")
}
