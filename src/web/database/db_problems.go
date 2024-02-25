package database

import (
	"database/sql"
	"errors"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/models"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// CreateProblemWithTxは，トランザクション内で新しい問題をデータベースに挿入する関数である．
// この関数は，与えられたトランザクションと問題モデルを使用して，問題のメタデータをProblemsテーブルに挿入する．
// 挿入操作が成功すると，挿入された行のIDが返される.
//
// この関数は，問題の作成に関連する全てのデータベース操作をトランザクションの一部として実行し，
// 操作中にエラーが発生した場合には，そのエラーを返す．エラーがなければ，新しく挿入された問題のIDを返す．
//
// パラメータ:
// - tx *sql.Tx: 実行中のトランザクション．
// - problem models.Problem: データベースに挿入する問題のモデル．
//
// 戻り値:
// - int: データベースに新しく挿入された問題のID．
// - error: 操作中に発生したエラー．成功時はnil．
func CreateProblemWithTx(tx *sql.Tx, problem models.Problem) (int, error) {
	var lastInsertId int64

	query := `INSERT INTO Problems (UserID, Title, Description, Difficulty) VALUES (?, ?, ?, ?)`
	result, execErr := tx.Exec(query, problem.UserID, problem.Title, problem.Description, problem.Difficulty)
	if execErr != nil {
		return 0, execErr // 直接エラーを返す
	}
	lastInsertId, execErr = result.LastInsertId() // 挿入された行のIDを取得
	if execErr != nil {
		return 0, execErr // 直接エラーを返す
	}

	return int(lastInsertId), nil // 挿入された行のIDを返す
}

// UpdateProblemは，指定されたIDの問題を更新する．
//
// この関数はデータベーストランザクションを用いて，問題の基本情報（Title, Description, Difficulty）の更新をアトミックに行うことを保証する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - problemID int: 更新対象の問題IDである．
// - problem models.Problem: 更新データを含む問題モデルである．
//
// 戻り値:
// - error: 更新操作に失敗した場合のエラー，または操作が成功した場合はnil．
//
// トランザクションを用いることで，更新プロセス中にエラーが発生した場合には，変更がロールバックされ，データベースの整合性を保つ．
func UpdateProblem(db *sql.DB, problemID int, problem models.Problem) error {
	err := WithTransaction(db, func(tx *sql.Tx) error {
		query := `UPDATE Problems SET Title = ?, Description = ?, Difficulty = ? WHERE ProblemID = ?`
		if _, err := tx.Exec(query, problem.Title, problem.Description, problem.Difficulty, problemID); err != nil {
			return err
		}
		return nil
	})

	// トランザクションエラー
	if err != nil {
		return err
	}

	return nil
}

// DeleteProblemは，指定された問題IDに関連する問題およびそれに紐付く全てのデータをデータベースから削除する．
// この処理には，問題自身のレコードの削除の他に，解答，テストケース結果など，問題に関連するデータの削除も含まれる．
// データベーストランザクションを使用して，削除操作がアトミックに行われることを保証する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - problemID int: 削除対象の問題IDである．
//
// 戻り値:
// - error: 削除操作に失敗した場合のエラー，または操作が成功した場合はnil．
//
// トランザクションを用いることで，更新プロセス中にエラーが発生した場合には，変更がロールバックされ，データベースの整合性を保つ．
func DeleteProblem(db *sql.DB, problemID int) error {
	err := WithTransaction(db, func(tx *sql.Tx) error {

		// NOTE: 問題に関連する全てのデータを安全に削除するために，関連データが存在する各テーブルに対して依存度の低いものから順番にDELETE文を実行する必要がある．
		if _, err := tx.Exec("DELETE FROM CaseResults WHERE SolutionID IN (SELECT SolutionID FROM Solutions WHERE ProblemID = ?)", problemID); err != nil {
			return err
		}

		if _, err := tx.Exec("DELETE FROM ResultDetails WHERE SolutionID IN (SELECT SolutionID FROM Solutions WHERE ProblemID = ?)", problemID); err != nil {
			return err
		}

		if _, err := tx.Exec("DELETE FROM Solutions WHERE ProblemID = ?", problemID); err != nil {
			return err
		}

		if _, err := tx.Exec("DELETE FROM Problems WHERE ProblemID = ?", problemID); err != nil {
			return err
		}
		return nil
	})

	// トランザクションエラー
	if err != nil {
		return err
	}

	return nil
}

// SelectProblemは，登録されている全問題のリストをデータベースから取得する．
// 各問題には，問題ID，ユーザーID，タイトル，説明文，難易度，作成日時，更新日時が含まれる．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
//
// 戻り値:
// - []models.Problem: 取得した問題のスライス．
// - error: データベース操作中にエラーが発生した場合の詳細．成功時はnil．
func SelectProblem(db *sql.DB) ([]models.Problem, error) {
	problems := []models.Problem{}

	// 問題の取得
	query := `SELECT ProblemID, UserID, Title, Description, Difficulty, CreatedAt, UpdatedAt FROM Problems`
	rows, err := db.Query(query)
	if err != nil {
		return nil, commonerrors.WrapDBError("SELECT", err)
	}
	defer rows.Close()

	problemMap := make(map[int]*models.Problem)
	for rows.Next() {
		var problem models.Problem
		if err := rows.Scan(&problem.ProblemID, &problem.UserID, &problem.Title, &problem.Description, &problem.Difficulty, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
			return nil, commonerrors.WrapDBError("ITERATING SELECTED SQL ROWS", err)
		}
		problemMap[problem.ProblemID] = &problem
		problems = append(problems, problem)
	}

	return problems, nil
}

// SelectProblemByUserIDは，特定のユーザーが作成した問題の詳細情報をデータベースから取得する．
// ユーザーIDに基づいて，そのユーザーが作成した全ての問題を取得する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - userID int: 問題の詳細情報を取得するユーザーのIDである．
//
// 戻り値:
// - []*models.Problem: 特定のユーザーが作成した問題のスライス．
// - error: データベース操作中にエラーが発生した場合の詳細．成功時はnil．
func SelectProblemByUserID(db *sql.DB, userID int) ([]*models.Problem, error) {
	problems := []*models.Problem{}

	// 問題の取得
	query := `SELECT ProblemID, UserID, Title, Description, Difficulty, CreatedAt, UpdatedAt FROM Problems WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, commonerrors.NewNotFoundError("Problem", "UserID", strconv.Itoa(userID))
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}
	defer rows.Close()

	problemMap := make(map[int]*models.Problem)
	for rows.Next() {
		var problem models.Problem
		if err := rows.Scan(&problem.ProblemID, &problem.UserID, &problem.Title, &problem.Description, &problem.Difficulty, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
			return nil, commonerrors.WrapDBError("ITERATING SELECTED SQL ROWS", err)
		}
		problemPtr := &problem
		problems = append(problems, problemPtr)
		problemMap[problem.ProblemID] = problemPtr
	}

	return problems, nil
}

// SelectProblemByProblemIDは，指定された問題IDに基づき，特定の問題の詳細情報をデータベースから取得する．
// 問題IDを指定して，その問題のID，ユーザーID，タイトル，説明，難易度，作成日時，更新日時を取得する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - problemID int: 詳細情報を取得する問題のIDである．
//
// 戻り値:
// - *models.Problem: 取得した問題の詳細情報を含むProblem型のポインタ．問題が存在しない場合はnil．
// - error: データベース操作中に発生したエラーの詳細．成功時はnil．
func SelectProblemByProblemID(db *sql.DB, problemID int) (*models.Problem, error) {
	var problem models.Problem

	// 問題の取得
	query := `SELECT ProblemID, UserID, Title, Description, Difficulty, CreatedAt, UpdatedAt FROM Problems WHERE ProblemID = ?`
	if err := db.QueryRow(query, problemID).Scan(&problem.ProblemID, &problem.UserID, &problem.Title, &problem.Description, &problem.Difficulty, &problem.CreatedAt, &problem.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 問題が見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("Problem", "ProblemID", strconv.Itoa(problemID))
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}

	return &problem, nil
}

// IsProblemOwnerは，指定されたユーザーが指定された問題の所有者かどうかを確認する．
// この関数は，指定されたユーザーIDと問題IDに基づいて，Problemsテーブルに存在するかどうかをチェックし，所有者であればnil，所有者でなければ関連するエラーを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - userID int: 所有者かどうかを確認したいユーザーのIDである．
// - problemID int: 所有者を確認したい問題のIDである．
//
// 戻り値:
// - error: データベース操作中に発生したエラーの詳細．成功時はnil．
func IsProblemOwner(db *sql.DB, userID int, problemID int) error {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM Problems WHERE ProblemID = ? AND UserID = ?)`
	if err := db.QueryRow(query, problemID, userID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 条件を満たす(ユーザID, 問題ID)が見つからないエラーを生成
			return commonerrors.NewAccessDeniedError("You do not have permission to access this resource")
			// return false, commonerrors.NewNotFoundError("Problem", "(UserID, ProblemID)", fmt.Sprintf("(%d, %d)", userID, problemID))
		}
		return commonerrors.WrapDBError("SELECT", err)
	}

	return nil
}
