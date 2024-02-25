package handlers

import (
	"database/sql"
	"net/http"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/minio"
	"procon_web_service/src/common/models"
	"procon_web_service/src/common/utils"
	"procon_web_service/src/web/database"
	webutils "procon_web_service/src/web/utils"
)

const (
	maxFileSize = 32 << 20 // 32MBのファイルサイズ制限
)

// UploadProblemHandlerは，新しい問題の投稿を処理するHTTPハンドラ関数である．
// この関数はHTTPリクエストから問題のメタデータと関連する入出力ファイルを解析し，それらをデータベースおよびMinIOに保存する．
// 問題のメタデータはリクエストボディから`models.Problem`構造体にデコードされ，入出力ファイルはマルチパートフォームデータとして処理される．
// この関数は認証情報の確認，マルチパートフォームデータのパース，ファイルの妥当性検証，問題メタデータとファイルの保存をトランザクション内で行う．
// 各ステップでエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
// 問題が正常に保存された場合，HTTPステータスコード201(Created)と保存された問題データをレスポンスとして返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 問題投稿処理を行う関数．
func UploadProblemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newProblem models.Problem

		// コンテキストから認証情報を取り出す
		if userClaims, ok := r.Context().Value("userClaims").(*models.Claims); !ok {
			// 認証情報が見つからない場合の処理
			return
		} else {
			newProblem.UserID = userClaims.UserID
		}

		if err := utils.ParseMultipartFormData(r, maxFileSize); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// ファイルのアップロードとMinIOへの保存に向けたフォーマットの確認
		if err := webutils.ValidateFiles(r.MultipartForm.File["input_file"], r.MultipartForm.File["output_file"]); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// JSON形式の問題メタデータを取得
		if err := utils.ParseProblemMetadata(r, &newProblem); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// トランザクションの開始
		tx, txErr := database.BeginTransaction(db)
		if txErr != nil {
			utils.SendErrorResponse(w, txErr)
			return
		}

		// [1] データベースに問題のメタデータを保存
		if problemID, err := database.CreateProblemWithTx(tx, newProblem); err != nil {
			tx.Rollback()
			utils.SendErrorResponse(w, err)
			return
		} else {
			newProblem.ProblemID = problemID // 割り振られた問題IDをProblem構造体にセット
		}

		// [2] 入力ファイルの保存
		if err := minio.UploadFileToMinIO(newProblem.ProblemID, r.MultipartForm.File["input_file"], "in"); err != nil {
			tx.Rollback()
			utils.SendErrorResponse(w, err)
			return
		}

		// [3] 出力ファイルの保存
		if err := minio.UploadFileToMinIO(newProblem.ProblemID, r.MultipartForm.File["output_file"], "out"); err != nil {
			tx.Rollback()
			utils.SendErrorResponse(w, err)
			return
		}

		// トランザクションのコミット( [1][2][3] が全て成功した時のみ)
		if err := tx.Commit(); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusCreated, newProblem)
	}
}

// UpdateProblemHandlerは，指定されたIDの問題を更新するHTTPハンドラ関数である．
// この関数はHTTPリクエストから問題の新しいメタデータと関連する入出力ファイルを解析し，それらをデータベースおよびMinIOに更新する．
// 問題のメタデータはリクエストボディから`models.Problem`構造体にデコードされ，入出力ファイルはマルチパートフォームデータとして処理される．
// この関数は認証情報の確認，マルチパートフォームデータのパース，ファイルの妥当性検証，既存の問題メタデータとファイルの更新を行う．
// まず，MinIOから古いファイルを削除し，新しいファイルを保存する．その後，データベースの問題メタデータを更新する．
// 各ステップでエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
// 問題が正常に更新された場合，HTTPステータスコード200(OK)と更新された問題データをレスポンスとして返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 問題更新処理を行う関数．
func UpdateProblemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var problem models.Problem
		// URLからProblemIDを取得
		if problemID, err := utils.GetIntVarFromRequest(r, "problem_id"); err != nil {
			utils.SendErrorResponse(w, err)
			return
		} else {
			// URLから修正する問題のProblemIDをセット
			problem.ProblemID = problemID
		}

		// コンテキストから認証情報を取り出す
		if userClaims, ok := r.Context().Value("userClaims").(*models.Claims); !ok {
			// 認証情報が見つからない場合の処理
			return
		} else {
			// 問題の作成者(=UserID)をJWTクレームからセット
			problem.UserID = userClaims.UserID
		}

		// リクエストのパース処理
		if err := utils.ParseMultipartFormData(r, maxFileSize); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// ファイルのアップロードとMinIOへの保存に向けたフォーマットの確認
		if err := webutils.ValidateFiles(r.MultipartForm.File["input_file"], r.MultipartForm.File["output_file"]); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// JSON形式の問題メタデータを取得
		if err := utils.ParseProblemMetadata(r, &problem); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// minIOの特定のバケットから古い問題の入出力データを削除(input/*, output/* まとめて)
		if err := minio.DeleteFileFromMinIO(minio.GetFileSaveName("", problem.ProblemID, "", "")); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// 入力ファイルの保存
		if err := minio.UploadFileToMinIO(problem.ProblemID, r.MultipartForm.File["input_file"], "input"); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// 出力ファイルの保存
		if err := minio.UploadFileToMinIO(problem.ProblemID, r.MultipartForm.File["output_file"], "output"); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// データベースに問題のメタデータを保存
		if err := database.UpdateProblem(db, problem.ProblemID, problem); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, problem)
	}
}

// DeleteProblemHandlerは，指定されたIDの問題とその関連データを削除するHTTPハンドラ関数である．
// この関数はURLパラメータから問題IDを抽出し，その問題に関連するデータベース内のメタデータとMinIO内のファイルを削除する．
// 削除処理は，まずMinIOから問題に関連する入出力ファイルを削除し，次にデータベースから問題のメタデータを削除することで行われる．
// この関数は，認証されたユーザーが自分の問題を削除することを許可するため，認証情報の確認も行う．
// 各ステップでエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
// 問題とその関連データが正常に削除された場合，HTTPステータスコード204(No Content)をレスポンスとして返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 問題削除処理を行う関数．
func DeleteProblemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからProblemIDを取得
		problemID, err := utils.GetIntVarFromRequest(r, "problem_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// minIOの特定のバケットから問題の入出力データを削除(input/*, output/* まとめて)
		if err := minio.DeleteFileFromMinIO(minio.GetFileSaveName("", problemID, "", "")); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// 特定の問題および関連したデータをデータベースから削除
		if err := database.DeleteProblem(db, problemID); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusNoContent, nil)
	}
}

// GetProblemHandlerは，登録されている全問題のリストを取得するHTTPハンドラ関数である．
// この関数はデータベースから全ての問題のメタデータを取得し，それらをレスポンスとして返す．
// 取得される問題のデータには，問題ID，作成者ID，タイトル，説明，難易度，作成日時，更新日時が含まれる．
// 問題データの取得に成功した場合，HTTPステータスコード200(OK)とともに問題リストをJSON形式で返す．
// データベース操作中にエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 問題リスト取得処理を行う関数．
func GetProblemHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		problems, err := database.SelectProblem(db)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, problems)
	}
}

// GetProblemByUserIDHandlerは，指定されたユーザーIDによって作成された問題のリストを取得するHTTPハンドラ関数である．
// この関数はURLパラメータからユーザーIDを取得し，そのIDに紐付く問題のメタデータをデータベースから取得する．
// 取得される問題のデータには，問題ID，作成者ID，タイトル，説明，難易度，作成日時，更新日時が含まれる．
// 問題データの取得に成功した場合，HTTPステータスコード200(OK)とともに問題リストをJSON形式で返す．
// ユーザーIDに紐付く問題が見つからない場合やデータベース操作中にエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 特定ユーザーの作成した問題リスト取得処理を行う関数．
func GetProblemByUserIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからUserIDを取得
		userID, err := utils.GetIntVarFromRequest(r, "user_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		problem, err := database.SelectProblemByUserID(db, userID)
		if err != nil {
			// ステータスコードの設定
			switch err.(type) {
			case *commonerrors.NotFoundError:
				// ユーザーが見つからない場合の処理
				utils.SendErrorResponse(w, err)
			case *commonerrors.ConstraintViolationError, *commonerrors.ForeignKeyViolationError, *commonerrors.DataMismatchError:
				// 制約エラー
				utils.SendErrorResponse(w, err)
			default:
				// その他のエラー
				utils.SendErrorResponse(w, err)
			}
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, problem)
	}
}

// GetProblemByProblemIDHandlerは，指定された問題IDの詳細情報を取得するHTTPハンドラ関数である．
// この関数はURLパラメータから問題IDを取得し，そのIDに紐付く問題のメタデータと関連するカテゴリーIDをデータベースから取得する．
// 取得される問題のデータには，問題ID，作成者ID，タイトル，説明，難易度，作成日時，更新日時，カテゴリーIDのリストが含まれる．
// 問題データの取得に成功した場合，HTTPステータスコード200(OK)とともに問題データをJSON形式で返す．
// 指定された問題IDの問題が見つからない場合やデータベース操作中にエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 特定の問題IDの詳細情報取得処理を行う関数．
func GetProblemByProblemIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからProblemIDを取得
		problemID, err := utils.GetIntVarFromRequest(r, "problem_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		problem, err := database.SelectProblemByProblemID(db, problemID)
		if err != nil {
			// 問題が見つからない場合の処理
			switch err.(type) {
			case *commonerrors.NotFoundError:
				// ユーザーが見つからない場合の処理
				utils.SendErrorResponse(w, err)
			case *commonerrors.ConstraintViolationError, *commonerrors.ForeignKeyViolationError, *commonerrors.DataMismatchError:
				// 制約エラー
				utils.SendErrorResponse(w, err)
			default:
				// その他のエラー
				utils.SendErrorResponse(w, err)
			}
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, problem)
	}
}
