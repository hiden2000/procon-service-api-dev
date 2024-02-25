package database

import (
	"database/sql"
	"errors"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/models"
	"strconv"
)

// CreateResultDetailは，ジャッジ結果をデータベースに保存する関数である．
// この関数は，解答IDとジャッジ結果の詳細を含むmodels.ResultDetail構造体を引数に取り，データベースに保存する．
// ジャッジ結果の詳細には，総テストケース数，正解数，不正解数，タイムリミット超過数，エラーメッセージが含まれる．
// また，各テストケースの結果もCaseResultsテーブルに保存される．
// この操作はデータベーストランザクション内で行われ，トランザクションが正常に完了しなかった場合はエラーが返される．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - solutionID int: ジャッジ結果が関連する解答のIDである．
// - resultDetail *models.ResultDetail: 保存するジャッジ結果の詳細である．
//
// 戻り値:
// - error: ジャッジ結果の保存に失敗した場合のエラー，または操作が成功した場合はnilである．
//
// トランザクションを用いることで，更新プロセス中にエラーが発生した場合には，変更がロールバックされ，データベースの整合性を保つ．
func CreateResultDetail(db *sql.DB, solutionID int, resultDetail *models.ResultDetail) error {
	err := WithTransaction(db, func(tx *sql.Tx) error {
		query := `INSERT INTO ResultDetails (SolutionID, TotalCases, CorrectCases, IncorrectCases, TimeLimitExceeded, ErrorMessage) VALUES (?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(query, solutionID, resultDetail.TotalCases, resultDetail.CorrectCases, resultDetail.IncorrectCases, resultDetail.TimeLimitExceeded, resultDetail.ErrorMessage)
		if err != nil {
			return err
		}

		for _, caseResult := range resultDetail.CaseResults {
			query = `INSERT INTO CaseResults (SolutionID, CaseName, Result, ExecutionTime) VALUES (?, ?, ?, ?)`
			_, err = tx.Exec(query, solutionID, caseResult.CaseName, caseResult.Result, caseResult.ExecutionTime)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// トランザクションエラー
	if err != nil {
		return err
	}

	return nil
}

// SelectResultDetailBySolutionIDは，特定の解答IDに対する判定結果を取得する関数である．
// この関数は，指定された解答IDに対応するResultDetailsテーブルからジャッジ結果の詳細を取得する．
// 取得したジャッジ結果には，総テストケース数，正解数，不正解数，タイムリミット超過数，エラーメッセージが含まれる．
// さらに，関連する各テストケースの結果もCaseResultsテーブルから取得され，models.ResultDetail構造体に格納される．
// 解答の詳細がデータベースに存在しない場合，NotFoundErrorが返される．
// この関数はデータベースからの情報の取得に失敗した場合にエラーを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタである．
// - solutionID int: 判定結果を取得したい解答のIDである．
//
// 戻り値:
// - *models.ResultDetail: 特定の解答IDに対する判定結果の詳細である．
// - error: 判定結果の取得に失敗した場合のエラー，または操作が成功した場合はnilである．
func SelectResultDetailBySolutionID(db *sql.DB, solutionID int) (*models.ResultDetail, error) {
	var resultDetail models.ResultDetail
	var caseResults []models.CaseResult

	err := db.QueryRow("SELECT TotalCases, CorrectCases, IncorrectCases, TimeLimitExceeded, ErrorMessage FROM ResultDetails WHERE SolutionID = ?", solutionID).Scan(&resultDetail.TotalCases, &resultDetail.CorrectCases, &resultDetail.IncorrectCases, &resultDetail.TimeLimitExceeded, &resultDetail.ErrorMessage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 解答の詳細が見つからないエラーを生成
			return nil, commonerrors.NewNotFoundError("ResultDetails", "SolutionID", strconv.Itoa(solutionID))
		}
		return nil, commonerrors.WrapDBError("SELECT", err)
	}

	rows, err := db.Query("SELECT CaseName, Result, ExecutionTime FROM CaseResults WHERE SolutionID = ?", solutionID)
	if errors.Is(err, sql.ErrNoRows) {
		return &resultDetail, nil
	} else if err != nil {
		return nil, commonerrors.WrapDBError("SELECT", err)
	}
	defer rows.Close()

	for rows.Next() {
		var caseResult models.CaseResult
		if err := rows.Scan(&caseResult.CaseName, &caseResult.Result, &caseResult.ExecutionTime); err != nil {
			return nil, commonerrors.WrapDBError("ITERATING SELECTED SQL ROWS", err)
		}
		caseResults = append(caseResults, caseResult)
	}
	resultDetail.CaseResults = caseResults

	return &resultDetail, nil
}
