package database

import (
	"database/sql"
	"errors"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/models"
	"strconv"
)

// CreateSolutionは，新しい解答をデータベースに保存する関数である．
// 引数として提供されたmodels.Solution構造体に基づき，解答情報をSolutionsテーブルに挿入する．
// 挿入操作はデータベーストランザクション内で行われ，挿入が成功すれば新しく生成された解答のIDを返し，失敗した場合はエラーを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - solution models.Solution: 保存する解答のデータを含む構造体．
//
// 戻り値:
// - int: 挿入された解答のID．
// - error: 操作が失敗した場合のエラー，またはnil．
//
// トランザクションを用いることで，更新プロセス中にエラーが発生した場合には，変更がロールバックされ，データベースの整合性を保つ．
func CreateSolution(db *sql.DB, solution models.Solution) (int, error) {
	var lastInsertId int64

	err := WithTransaction(db, func(tx *sql.Tx) error {
		query := `INSERT INTO Solutions (UserID, ProblemID, LanguageID, Code) VALUES (?, ?, ?, ?)`
		result, execErr := tx.Exec(query, solution.UserID, solution.ProblemID, solution.LanguageID, solution.Code)
		if execErr != nil {
			return execErr
		}
		lastInsertId, execErr = result.LastInsertId() // 挿入された行のIDを取得
		return execErr
	})

	// トランザクションエラー
	if err != nil {
		return 0, err // エラーがあればIDの代わりに0を返す
	}

	return int(lastInsertId), nil // 挿入された行のIDを返す
}

// SelectSolutionByUserIDは，特定ユーザーの全解答をデータベースから取得する関数である．
// 指定されたユーザーIDに基づき，Solutionsテーブルから該当する解答のリストを取得し，models.Solution構造体のスライスとして返す．
// 解答が一つも見つからない場合は，NotFoundErrorを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - userID int: 解答を取得したいユーザーのID．
//
// 戻り値:
// - []models.Solution: 特定ユーザーの解答リスト．
// - error: 操作が失敗した場合のエラー，またはnil．
func SelectSolutionByUserID(db *sql.DB, userID int) ([]models.Solution, error) {
	solutions := []models.Solution{}

	query := `SELECT SolutionID, UserID, ProblemID, LanguageID, Code, SubmittedAt FROM Solutions WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 解答が見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("Solution", "UserID", strconv.Itoa(userID))
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}
	defer rows.Close()

	for rows.Next() {
		var solution models.Solution
		if err := rows.Scan(&solution.SolutionID, &solution.UserID, &solution.ProblemID, &solution.LanguageID, &solution.Code, &solution.SubmittedAt); err != nil {
			return nil, commonerrors.WrapDBError("ITERATING SELECTED SQL ROWS", err)
		}
		solutions = append(solutions, solution)
	}

	return solutions, nil
}

// SelectSolutionByProblemIDは，特定問題に関する全解答をデータベースから取得する関数である．
// 指定された問題IDに基づき，Solutionsテーブルから該当する解答のリストを取得し，models.Solution構造体のスライスとして返す．
// 解答が一つも見つからない場合は，NotFoundErrorを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - problemID int: 解答を取得したい問題のID．
//
// 戻り値:
// - []models.Solution: 特定問題の解答リスト．
// - error: 操作が失敗した場合のエラー，またはnil．
func SelectSolutionByProblemID(db *sql.DB, problemID int) ([]models.Solution, error) {
	solutions := []models.Solution{}

	query := `SELECT SolutionID, UserID, ProblemID, LanguageID, Code, SubmittedAt FROM Solutions WHERE ProblemID = ?`
	rows, err := db.Query(query, problemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 解答が見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("Solution", "ProblemID", strconv.Itoa(problemID))
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}
	defer rows.Close()

	for rows.Next() {
		var solution models.Solution
		if err := rows.Scan(&solution.SolutionID, &solution.UserID, &solution.ProblemID, &solution.LanguageID, &solution.Code, &solution.SubmittedAt); err != nil {
			return nil, commonerrors.WrapDBError("ITERATING SELECTED SQL ROWS", err)
		}
		solutions = append(solutions, solution)
	}

	return solutions, nil
}

// SelectSolutionBySolutionIDは，特定の解答IDに対する詳細情報をデータベースから取得する関数である．
// 指定された解答IDに基づき，Solutionsテーブルから該当する解答のデータを取得し，models.Solution構造体として返す．
// 解答が見つからない場合は，NotFoundErrorを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
// - solutionID int: 詳細情報を取得したい解答のID．
//
// 戻り値:
// - *models.Solution: 特定の解答の詳細情報．
// - error: 操作が失敗した場合のエラー，またはnil．
func SelectSolutionBySolutionID(db *sql.DB, solutionID int) (*models.Solution, error) {
	var solution models.Solution

	query := `SELECT SolutionID, UserID, ProblemID, LanguageID, Code, SubmittedAt FROM Solutions WHERE SolutionID = ?`
	if err := db.QueryRow(query, solutionID).Scan(&solution.SolutionID, &solution.UserID, &solution.ProblemID, &solution.LanguageID, &solution.Code, &solution.SubmittedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 解答が見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("Solution", "SolutionID", strconv.Itoa(solutionID))
		}
		commonerrors.WrapDBError("SELECT", err)
	}

	return &solution, nil
}
