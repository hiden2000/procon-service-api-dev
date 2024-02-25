package database

import (
	"database/sql"
	"fmt"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/web/config"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectToDBは，アプリケーションのデータベース接続を初期化し，データベース接続へのポインタを返す関数である．
// この関数は，環境変数から読み込まれたデータベース設定を使用して，MySQLデータベースへの接続文字列を構築し，
// sql.Openを呼び出してデータベース接続を確立する．接続が成功すれば，*sql.DB型のデータベース接続へのポインタとnilを返し，失敗した場合はnilとエラーを返す．
//
// 戻り値:
// - *sql.DB: データベース接続へのポインタ．
// - error: 接続に失敗した場合のエラー，またはnil．
func ConnectToDB() (*sql.DB, error) {

	dbconfig := config.NewDatabaseConfig()
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbconfig.User, dbconfig.Password, dbconfig.Host, dbconfig.Database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, &commonerrors.DBError{Op: "mysql", Err: err}
	}

	return db, nil
}

// WithTransactionは，データベーストランザクションを作成し，提供された関数fnをそのトランザクション内で実行する関数である．
// この関数は，トランザクションを開始するためにdb.Beginを呼び出し，提供されたfn関数にトランザクションを渡す．
// fn関数がエラーなしで完了した場合，トランザクションはコミットされる．fn関数実行中にエラーが発生した場合，
// トランザクションはロールバックされ，エラーが返される．コミットまたはロールバックのいずれかの操作でエラーが発生した場合は，そのエラーを含むTransactionErrorが返される．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - fn func(*sql.Tx) error: トランザクション内で実行される関数．トランザクションオブジェクトを引数とし，エラーを返すことができる．
//
// 戻り値:
// - *commonerrors.TransactionError: トランザクション操作中に発生したエラーをラップしたカスタムエラー，またはnil．
func WithTransaction(db *sql.DB, fn func(*sql.Tx) error) *commonerrors.TransactionError {
	tx, err := db.Begin()
	if err != nil {
		return commonerrors.NewTransactionError("begin", err)
	}
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return commonerrors.NewTransactionError("rollback", rbErr)
		}
		return commonerrors.NewTransactionError("execute", err)
	}
	if err := tx.Commit(); err != nil {
		return commonerrors.NewTransactionError("commit", err)
	}
	return nil
}

// BeginTransactionは，データベースに対して新しいトランザクションを開始する．
// 成功した場合，開始されたトランザクションとnilエラーを返す．エラーが発生した場合は，トランザクションをnilとして，エラーを返す．
//
// パラメータ:
// - db *sql.DB: トランザクションを開始するデータベースのインスタンス．
//
// 戻り値:
// - *sql.Tx: 開始されたトランザクション．
// - *commonerrors.TransactionError: トランザクション開始時に発生したエラー．
func BeginTransaction(db *sql.DB) (*sql.Tx, *commonerrors.TransactionError) {
	tx, err := db.Begin()
	if err != nil {
		return nil, commonerrors.NewTransactionError("begin", err)
	}
	return tx, nil
}

// ExecuteInTransactionは，与えられたトランザクション内で特定の操作を実行する．
// 操作が成功した場合はnilを，失敗した場合はエラーを返す．
//
// パラメータ:
// - tx *sql.Tx: 操作を実行するトランザクション．
// - fn func(*sql.Tx) error: トランザクション内で実行する操作を定義する関数．
//
// 戻り値:
// - *commonerrors.TransactionError: 操作実行時に発生したエラー．
func ExecuteInTransaction(tx *sql.Tx, fn func(*sql.Tx) error) *commonerrors.TransactionError {
	if err := fn(tx); err != nil {
		return commonerrors.NewTransactionError("execute", err)
	}
	return nil
}

// EndTransactionは，トランザクションをコミットまたはロールバックし，トランザクションを終了する．
// エラーが引数として渡された場合，トランザクションはロールバックされる．エラーがnilの場合，トランザクションはコミットされる．
//
// パラメータ:
// - tx *sql.Tx: 終了するトランザクション．
// - err *commonerrors.TransactionError: トランザクション実行中に発生したエラー．
//
// 戻り値:
// - *commonerrors.TransactionError: コミットまたはロールバック時に発生したエラー．
func EndTransaction(tx *sql.Tx, err *commonerrors.TransactionError) *commonerrors.TransactionError {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return commonerrors.NewTransactionError("rollback", rbErr)
		}
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return commonerrors.NewTransactionError("commit", commitErr)
	}
	return nil
}
