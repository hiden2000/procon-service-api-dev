package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/models"
	"procon_web_service/src/common/utils"
	"procon_web_service/src/web/database"
	webutils "procon_web_service/src/web/utils"
)

// RegisterUserHandlerは，新規ユーザー登録を処理するHTTPハンドラ関数である．
// この関数は，HTTPリクエストからユーザー情報をデコードし，デコードされたユーザー情報をデータベースに保存する．
// ユーザーのパスワードはハッシュ化された後，データベースに保存される．
// さらに，ユーザー情報からJWTトークンが生成され，クライアントに返される．
// データベース操作とJWTトークンの生成は，トランザクション内で実行される．
// これにより，どちらかの処理でエラーが発生した場合，トランザクションがロールバックされ，データベースへの変更がなかったことになる．
// 処理が成功した場合，HTTPステータスコード201(Created)と生成されたJWTトークンがクライアントに返される．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: 新規ユーザー登録処理を行う関数．
func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.User

		if err := utils.DecodeRequestBody(r, &newUser); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// パスワードのハッシュ化
		if hashedPassword, err := webutils.GenerateHashedPassword(newUser.Password); err != nil {
			utils.SendErrorResponse(w, err)
			return
		} else {
			newUser.Password = hashedPassword
		}

		// トランザクションの開始
		tx, txErr := database.BeginTransaction(db)
		if txErr != nil {
			utils.SendErrorResponse(w, commonerrors.NewTransactionError("begin", txErr))
			return
		}

		// [1] ユーザー情報をデータベースに保存
		if userID, err := database.CreateUserWithTx(tx, newUser); err != nil {
			tx.Rollback()
			utils.SendErrorResponse(w, err)
			return
		} else {
			newUser.UserID = userID
		}

		// [2] JWTトークンの生成
		token, err := webutils.GenerateJWT(newUser)
		if err != nil {
			tx.Rollback()
			utils.SendErrorResponse(w, err)
			return
		}

		// トランザクションのコミット( [1][2] が両方成功した時のみ)
		if err := tx.Commit(); err != nil {
			utils.SendErrorResponse(w, commonerrors.NewTransactionError("commit", err))
			return
		}

		utils.SendJSONResponse(w, http.StatusCreated, map[string]interface{}{
			"token": token,
			"user":  newUser,
		})
	}
}

// UpdateUserHandlerは，認証されたユーザーのプロファイル情報を更新するHTTPハンドラ関数である．
// この関数は，リクエストボディからユーザープロファイル情報をデコードし，その情報を使用してデータベース内の対応するユーザー情報を更新する．
// 更新操作は，認証情報に基づいたユーザーIDで特定されたユーザーに対してのみ行われる．
// 更新が成功した場合，HTTPステータスコード200(OK)と更新されたユーザープロファイル情報を含むレスポンスが返される．
// リクエストボディの解析に失敗した場合やデータベースの更新操作中にエラーが発生した場合は，適切なHTTPステータスコードとエラーメッセージで応答する．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: ユーザープロファイル更新処理を行う関数．
func UpdateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// コンテキストから認証情報を取り出す
		userClaims, ok := r.Context().Value("userClaims").(*models.Claims)
		if !ok {
			// 認証情報が見つからない場合の処理
			return
		}

		var userProfile models.UserProfile
		if err := utils.DecodeRequestBody(r, &userProfile); err != nil {
			utils.SendErrorResponse(w, err)
			return
		} else {
			userProfile.UserID = userClaims.UserID // ユーザープロファイルにuserIDをセット
		}

		if err := database.UpdateUser(db, userProfile.UserID, userProfile); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, userProfile)
	}
}

// LoginUserHandlerは，ユーザーログインを処理するHTTPハンドラ関数である．
// この関数はHTTPリクエストからユーザーの認証情報をデコードし，デコードされた認証情報に基づいてユーザー認証を行う．
// 認証に成功した場合，認証されたユーザーに対してJWTトークンが生成され，クライアントに返される．
// 認証に失敗した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
// この関数では，ユーザー名がデータベースに存在するかどうかを確認し，存在する場合はパスワードが一致するかどうかを検証する．
// パスワードが一致した場合，ユーザー認証が成功したとみなし，JWTトークンが生成される．
// 処理が成功した場合，HTTPステータスコード200(OK)と生成されたJWTトークン，およびユーザープロファイルがクライアントに返される．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: ユーザーログイン処理を行う関数．
func LoginUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials models.UserCredentials

		if err := utils.DecodeRequestBody(r, &credentials); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// ユーザー認証
		user, err := database.SelectUserByUsername(db, credentials.Username)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		if err := webutils.VerifyPassword(user.Password, credentials.Password); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// 認証成功後のJWTトークン生成
		token, err := webutils.GenerateJWT(*user)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		response := struct {
			Token string             `json:"token"`
			User  models.UserProfile `json:"user"`
		}{
			Token: token,
			User: models.UserProfile{
				UserID:   user.UserID,
				Username: user.Username,
			},
		}

		utils.SendJSONResponse(w, http.StatusOK, response)
	}
}

// LogoutUserHandlerは，ユーザーログアウトを処理するHTTPハンドラ関数である．
// この関数はHTTPリクエストヘッダーからJWTトークンを抽出し，そのトークンをRedisから削除してセッションを無効化することでユーザーログアウトを実現する．
// トークンの削除に成功した場合，HTTPステータスコード200(OK)と空のレスポンスボディがクライアントに返される．
// トークンの削除に失敗した場合，適切なHTTPステータスコードとエラーメッセージで応答する．
// この関数は，ユーザーがログイン状態であることを前提としており，ログアウト処理を通じてユーザーのセッションを終了させる．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタは不要だが，一貫性のために引数に含まれている．
//
// 戻り値:
// - http.HandlerFunc: ユーザーログアウト処理を行う関数．
func LogoutUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redisからトークンを削除してセッションを無効化
		if err := webutils.EraseToken(webutils.ExtractToken(r)); err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, nil)
	}
}

// GetUserByUserIDHandlerは，指定されたユーザーIDに基づいてユーザープロファイル情報を取得するHTTPハンドラ関数である．
// この関数はURLパラメータからユーザーIDを取得し，そのIDを使用してデータベースからユーザー情報を検索する．
// ユーザー情報が見つかった場合，HTTPステータスコード200(OK)とともにユーザープロファイル情報を含むレスポンスボディが返される．
// ユーザー情報が見つからなかった場合や，データベース検索時にエラーが発生した場合は，適切なHTTPステータスコードとエラーメッセージで応答する．
// ユーザープロファイル情報にはユーザーIDとユーザー名が含まれる．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: ユーザープロファイル取得処理を行う関数．
func GetUserByUserIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLからUserIDを取得
		userID, err := utils.GetIntVarFromRequest(r, "user_id")
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		user, err := database.SelectUserByUserID(db, userID)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		userProfile := models.UserProfile{
			UserID:   user.UserID,
			Username: user.Username,
		}

		utils.SendJSONResponse(w, http.StatusOK, userProfile)
	}
}

// GetUserByUsernameHandlerは，指定されたユーザー名に基づいてユーザープロファイル情報を取得するHTTPハンドラ関数である．
// この関数はクエリパラメータからユーザー名を取得し，その名前を使用してデータベースからユーザー情報を検索する．
// ユーザー情報が見つかった場合，HTTPステータスコード200(OK)とともにユーザープロファイル情報を含むレスポンスボディが返される．
// ユーザー名がクエリパラメータに存在しない，ユーザー情報が見つからなかった場合や，データベース検索時にエラーが発生した場合は，適切なHTTPステータスコードとエラーメッセージで応答する．
// ユーザープロファイル情報にはユーザーIDとユーザー名が含まれる．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - http.HandlerFunc: ユーザープロファイル取得処理を行う関数．
func GetUserByUsernameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータからUsernameを取得
		userName := r.URL.Query().Get("username")
		if userName == "" {
			utils.SendErrorResponse(w, commonerrors.WrapRequestVariableError("username", errors.New("username query parameter is required")))
			return
		}

		user, err := database.SelectUserByUsername(db, userName)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		userProfile := models.UserProfile{
			UserID:   user.UserID,
			Username: user.Username,
		}

		utils.SendJSONResponse(w, http.StatusOK, userProfile)
	}
}
