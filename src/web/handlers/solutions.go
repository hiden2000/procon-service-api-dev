package handlers

import (
	"database/sql"
	"net/http"
	"procon_web_service/src/common/models"
	"procon_web_service/src/common/utils"
	"procon_web_service/src/web/database"
)

// SubmitSolutionHandlerは，ユーザーからの解答提出を処理するHTTPハンドラ関数である．
// この関数はHTTPリクエストから解答データをデコードし，デコードされた解答をデータベースに保存する．
// 解答データは，リクエストボディから`models.Solution`構造体にデコードされ，URLパラメータから問題IDを取得して構造体にセットする．
// ユーザーの認証情報は，HTTPリクエストのコンテキストから取得し，解答にユーザーIDをセットする．
// 解答のデータベースへの保存が成功した場合，HTTPステータスコード201(Created)と保存された解答データをレスポンスとして返す．
// 各ステップでエラーが発生した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
// 非同期処理のトリガーはWebSocketHandler内で行われるため，この関数ではデータベースへの保存と初期応答のみを担当する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 解答提出処理を行う関数．
func SubmitSolutionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var solution models.Solution

		// コンテキストから認証情報を取り出す
		if userClaims, ok := r.Context().Value("userClaims").(*models.Claims); !ok {
			// 認証情報が見つからない場合の処理
			return
		} else {
			solution.UserID = userClaims.UserID // 問題の作成者(=UserID)をJWTクレームからsolution構造体にセット
		}

		// URLからProblemIDを取得
		if problemID, err := utils.GetIntVarFromRequest(r, "problem_id"); err != nil {
			utils.SendErrorResponse(w, err)
			return
		} else {
			solution.ProblemID = problemID
		}

		if err := utils.DecodeRequestBody(r, &solution); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		if solutionID, err := database.CreateSolution(db, solution); err != nil {
			utils.SendErrorResponse(w, err)
			return
		} else {
			solution.SolutionID = solutionID // データベースから割り振られた解答番号をセット
		}

		utils.SendJSONResponse(w, http.StatusCreated, solution)
	}
}

// GetSolutionDetailsHandlerは，特定の解答IDに対する詳細情報を取得するHTTPハンドラ関数である．
// この関数はURLパラメータから解答IDを取得し，対応する解答の詳細情報をデータベースから検索する．
// 取得した解答データは，`models.Solution`構造体に格納され，HTTPレスポンスとしてクライアントに返される．
// 解答IDに対応する解答がデータベースに存在しない場合や，データベースからのデータ取得に失敗した場合には，適切なエラーメッセージとHTTPステータスコードで応答する．
// 成功した場合，HTTPステータスコード200(OK)と共に解答の詳細情報を含むレスポンスを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 特定の解答の詳細情報を取得する処理を行う関数．
func GetSolutionDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからSolutionIDを取得
		solutionID, err := utils.GetIntVarFromRequest(r, "solution_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		solution, err := database.SelectSolutionBySolutionID(db, solutionID)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, solution)
	}
}

// GetSolutionsByUserIDHandlerは，指定されたユーザIDに関連する提出コード一覧を取得するHTTPハンドラ関数である．
// この関数は，リクエストからユーザIDを取得し，そのユーザIDに紐づく提出コードの一覧をデータベースから検索する．
// 提出コードは，`models.Solution`構造体のスライスとしてクライアントに返される．
// データベースからの検索に失敗した場合や，該当する提出コードが存在しない場合には，適切なエラーメッセージと共にエラーレスポンスを返す．
// 検索が成功した場合は，HTTPステータスコード200(OK)と共に，提出コードの一覧を含むレスポンスを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 指定されたユーザIDの提出コード一覧を取得する処理を行う関数．
func GetSolutionsByUserIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからUserIDを取得
		userID, err := utils.GetIntVarFromRequest(r, "user_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		solutions, err := database.SelectSolutionByUserID(db, userID)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, solutions)
	}
}

// GetSolutionsByProblemIDHandlerは，指定された問題IDに関連する提出コード一覧を取得するHTTPハンドラ関数である．
// この関数は，リクエストから問題IDを取得し，その問題IDに紐づく提出コードの一覧をデータベースから検索する．
// 提出コードは，`models.Solution`構造体のスライスとしてクライアントに返される．
// データベースからの検索に失敗した場合や，該当する提出コードが存在しない場合には，適切なエラーメッセージと共にエラーレスポンスを返す．
// 検索が成功した場合は，HTTPステータスコード200(OK)と共に，提出コードの一覧を含むレスポンスを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 指定された問題IDの提出コード一覧を取得する処理を行う関数．
func GetSolutionsByProblemIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからUserIDを取得
		problemID, err := utils.GetIntVarFromRequest(r, "problem_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		solutions, err := database.SelectSolutionByProblemID(db, problemID)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, solutions)
	}
}

// GetSolutionResultHandlerは，特定の解答IDに対する判定結果を取得するHTTPハンドラ関数である．
// この関数は，リクエストから解答IDを取得し，その解答IDに紐づく判定結果をデータベースから検索する．
// 判定結果は，`models.ResultDetail`構造体でクライアントに返される．
// データベースからの検索に失敗した場合や，該当する判定結果が存在しない場合には，適切なエラーメッセージと共にエラーレスポンスを返す．
// 検索が成功した場合は，HTTPステータスコード200(OK)と共に，判定結果を含むレスポンスを返す．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 特定の解答の判定結果を取得する処理を行う関数．
func GetSolutionResultHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからSolutionIDを取得
		solutionID, err := utils.GetIntVarFromRequest(r, "solution_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		resultDetail, err := database.SelectResultDetailBySolutionID(db, solutionID)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, resultDetail)
	}
}
