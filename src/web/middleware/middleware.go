package middleware

import (
	"context"
	"database/sql"
	"net/http"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/utils"
	"procon_web_service/src/web/database"
	webutils "procon_web_service/src/web/utils"
)

// AuthMiddlewareは，JWTトークンに基づいた認証機能を提供するミドルウェアである．
// このミドルウェアは，受け取ったHTTPリクエストに含まれるJWTトークンを検証し，トークンが有効である場合にのみ後続のハンドラへのリクエスト処理を許可する．
// トークンの検証には，utils.IsUserAuthenticated関数が使用され，認証に成功した場合，認証されたユーザーのクレーム情報がリクエストのコンテキストに"userClaims"として追加される．
// 認証に失敗した場合，HTTPステータスコード401(Unauthorized)とともにエラーメッセージがクライアントに送信される．
// 認証に成功した場合のみ，コンテキストを更新したリクエストで次のハンドラが呼び出される．この仕組みにより，認証が必要なAPIエンドポイントのセキュリティが強化される．
//
// パラメータ:
// - next http.Handler: 認証に成功した場合に実行される次のハンドラ関数．
//
// 戻り値:
// - http.Handler: 認証機能を追加したミドルウェア関数．
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := webutils.IsUserAuthenticated(r)
		if err != nil {
			utils.SendErrorResponse(w, err)
			return
		}

		// 認証に成功した場合，認証情報をリクエストのコンテキストに追加
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		// コンテキストを更新したリクエストで次のハンドラを呼び出す
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ProblemOwnerMiddlewareFactoryは，特定の問題の所有者であるかどうかを確認するミドルウェアを生成するファクトリ関数である．
// 生成されるミドルウェアは，HTTPリクエストから問題IDを抽出し，リクエストを行ったユーザーがその問題の所有者であるかをデータベースで確認する．
// ユーザー認証は，JWTトークンに基づいて行われる．
// 問題の所有者でない場合，HTTPステータスコード403(Forbidden)とエラーメッセージがクライアントに送信される．問題の所有者である場合のみ，次のハンドラが呼び出される．
// このミドルウェアは，問題の編集や削除など，特定の問題に対する操作を行うAPIエンドポイントにおいて，操作が問題の所有者によってのみ行われることを保証するために使用される．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - func(http.HandlerFunc) http.HandlerFunc: 指定されたhttp.HandlerFuncに対して問題の所有者確認機能を追加するミドルウェアを生成する関数．
func ProblemOwnerMiddlewareFactory(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// ユーザー照合のためJWTクレームの認証
			claims, err := webutils.IsUserAuthenticated(r)
			if err != nil {
				utils.SendErrorResponse(w, err)
				return
			}

			// URLからProblemIDを取得
			problemID, err := utils.GetIntVarFromRequest(r, "problem_id")
			if err != nil {
				utils.SendErrorResponse(w, err)
				return
			}

			// 問題の所有者か確認
			if err := database.IsProblemOwner(db, claims.UserID, problemID); err != nil {
				utils.SendErrorResponse(w, err)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// UserAuthMiddlewareFactoryは，特定のユーザー関連操作が認証されたユーザー自身によってのみ行われることを保証するミドルウェアを生成するファクトリ関数である．
// 生成されるミドルウェアは，HTTPリクエストからユーザーIDを抽出し，リクエストを行ったユーザーが操作しようとしているリソースの所有者であるかを確認する．
// ユーザー認証はJWTトークンに基づいて行われ，認証されたユーザーのクレーム情報とリクエストURLのユーザーIDが一致することが確認される．
// ユーザーIDが一致しない場合，HTTPステータスコード403(Forbidden)とエラーメッセージがクライアントに送信される．一致する場合のみ，次のハンドラが呼び出される．
// このミドルウェアは，ユーザープロファイルの更新や特定のユーザーに紐づくデータの取得など，ユーザー自身に関連する操作を行うAPIエンドポイントにおいて使用される．
//
// パラメータ:
// - db *sql.DB: データベース接続へのポインタ．
//
// 戻り値:
// - func(http.HandlerFunc) http.HandlerFunc: 指定されたhttp.HandlerFuncに対してユーザー認証機能を追加するミドルウェアを生成する関数．
func UserAuthMiddlewareFactory(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// ユーザー照合のためJWTクレームの認証
			claims, err := webutils.IsUserAuthenticated(r)
			if err != nil {
				utils.SendErrorResponse(w, err)
				return
			}

			// URLからUserIDを取得
			userID, err := utils.GetIntVarFromRequest(r, "user_id")
			if err != nil {
				utils.SendErrorResponse(w, err)
				return
			}

			// トークンuserIDとURLのuserIDが一致するか確認
			if userID != claims.UserID {
				utils.SendErrorResponse(w, commonerrors.NewAccessDeniedError("You do not have permission to access this resource"))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
