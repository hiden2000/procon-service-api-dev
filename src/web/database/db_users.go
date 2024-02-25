package database

import (
	"database/sql"
	"errors"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/models"
	"strconv"
)

// CreateUserWithTxは，新しいユーザー情報をデータベースに保存する関数である．
// 引数として提供されたmodels.User構造体に基づき，ユーザー情報をUsersテーブルに挿入する．
// 挿入操作はデータベーストランザクション内で行われ，挿入が成功すれば新しく生成されたユーザーのIDを返し，失敗した場合はエラーを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - user models.User: 保存するユーザーのデータを含む構造体．
//
// 戻り値:
// - int: 挿入されたユーザーのID．
// - error: 操作が失敗した場合のエラー，またはnil．
func CreateUserWithTx(tx *sql.Tx, user models.User) (int, error) {
	var lastInsertId int64
	query := `INSERT INTO Users (Username, Password) VALUES (?, ?)`
	result, execErr := tx.Exec(query, user.Username, user.Password)
	if execErr != nil {
		return 0, commonerrors.WrapDBError("insert", execErr) // 直接エラーを返す
	}
	lastInsertId, execErr = result.LastInsertId() // 挿入された行のIDを取得
	if execErr != nil {
		return 0, commonerrors.WrapDBError("insert", execErr) // 直接エラーを返す
	}

	return int(lastInsertId), nil // 挿入された行のIDを返す
}

// UpdateUserは，新しいユーザー情報をデータベースに保存する関数である．
// 引数として提供されたmodels.User構造体に基づき，ユーザー情報をUsersテーブルに挿入する．
// 挿入操作はデータベーストランザクション内で行われ，挿入が成功すれば新しく生成されたユーザーのIDを返し，失敗した場合はエラーを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - user models.User: 保存するユーザーのデータを含む構造体．
//
// 戻り値:
// - int: 挿入されたユーザーのID．
// - error: 操作が失敗した場合のエラー，またはnil．
func UpdateUser(db *sql.DB, userID int, profile models.UserProfile) error {
	err := WithTransaction(db, func(tx *sql.Tx) error {
		query := `UPDATE Users SET Username = ? WHERE UserID = ?`
		_, err := tx.Exec(query, profile.Username, userID)
		return err
	})

	// トランザクションエラー
	if err != nil {
		return err
	}

	return nil
}

// SelectUserByUsernameは，ユーザー名に基づいてユーザー情報を取得する関数である．
// 指定されたユーザー名に一致するユーザー情報をUsersテーブルから検索し，models.User構造体として返す．
// ユーザーが見つからない場合は，NotFoundErrorを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - username string: 情報を取得したいユーザーのユーザー名．
//
// 戻り値:
// - *models.User: 検索されたユーザーの情報．
// - error: 操作が失敗した場合のエラー，またはnil．
func SelectUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User

	query := `SELECT UserID, Username, Password FROM Users WHERE Username = ?`
	if err := db.QueryRow(query, username).Scan(&user.UserID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ユーザーが見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("User", "Username", username)
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}

	return &user, nil
}

// SelectUserByIDは，ユーザーIDによるユーザー情報を取得する関数である．
// 指定されたユーザーIDに一致するユーザー情報をUsersテーブルから検索し，models.User構造体として返す．
// ユーザーが見つからない場合は，NotFoundErrorを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - userID int: 情報を取得したいユーザーのID．
//
// 戻り値:
// - *models.User: 検索されたユーザーの情報．
// - error: 操作が失敗した場合のエラー，またはnil．
func SelectUserByUserID(db *sql.DB, userID int) (*models.User, error) {
	var user models.User

	query := `SELECT UserID, Username, Password FROM Users WHERE UserID = ?`
	if err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ユーザーが見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("User", "UserID", strconv.Itoa(userID))
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}

	return &user, nil
}

// CheckUserExistsは，指定されたユーザー名が既に存在するかをチェックする関数である．
// Usersテーブルを検索し，指定されたユーザー名が存在するかどうかの真偽値を返す．
// ユーザー名の存在チェックはSELECT EXISTSを用いて行われる．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - username string: 存在を確認したいユーザー名．
//
// 戻り値:
// - bool: 指定されたユーザー名が存在する場合はtrue，そうでない場合はfalse．
// - error: 操作が失敗した場合のエラー，またはnil．
func CheckUserExists(db *sql.DB, username string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM Users WHERE Username = ?)`
	if err := db.QueryRow(query, username).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ユーザーが見つからないエラーを生成
			return false, commonerrors.NewNotFoundError("User", "Username", username)
		}
		return false, commonerrors.WrapDBError("SELECT", err)
	}

	return exists, nil
}
