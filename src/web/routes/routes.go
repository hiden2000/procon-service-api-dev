package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	cmiddleware "procon_web_service/src/common/middleware"
	"procon_web_service/src/web/handlers"
	"procon_web_service/src/web/middleware"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(router *mux.Router, db *sql.DB) {
	router.Use(cmiddleware.LoggingMiddleware)

	publicRoutes := router.PathPrefix("/api").Subrouter()

	// 1. 認証の不要なAPIルート
	// 問題に関するAPI
	publicRoutes.HandleFunc("/problems", handlers.GetProblemHandler(db)).Methods(http.MethodGet)                         // 全ての問題の取得
	publicRoutes.HandleFunc("/problems/{problem_id}", handlers.GetProblemByProblemIDHandler(db)).Methods(http.MethodGet) // 指定された問題IDの問題概要を取得
	publicRoutes.HandleFunc("/users/{user_id}/problems", handlers.GetProblemByUserIDHandler(db)).Methods(http.MethodGet) // ユーザーIDに基づく問題の取得

	// ユーザーに関するAPI
	publicRoutes.HandleFunc("/users", handlers.RegisterUserHandler(db)).Methods(http.MethodPost)             // ユーザー登録
	publicRoutes.HandleFunc("/users/login", handlers.LoginUserHandler(db)).Methods(http.MethodPost)          // ログイン
	publicRoutes.HandleFunc("/users/{user_id}", handlers.GetUserByUserIDHandler(db)).Methods(http.MethodGet) // 指定されたユーザーIDのユーザーを取得
	publicRoutes.HandleFunc("/users", handlers.GetUserByUsernameHandler(db)).Methods(http.MethodGet)         // ユーザー名に基づくユーザープロファイルの取得(クエリパラメータ?usernameで取得)

	// 解答に関するAPI
	publicRoutes.HandleFunc("/solutions/{solution_id}", handlers.GetSolutionDetailsHandler(db)).Methods(http.MethodGet)              // 指定された解答IDの解答を取得
	publicRoutes.HandleFunc("/solutions/{solution_id}/result", handlers.GetSolutionResultHandler(db)).Methods(http.MethodGet)        // 解答の結果の取得
	publicRoutes.HandleFunc("/problems/{problem_id}/solutions", handlers.GetSolutionsByProblemIDHandler(db)).Methods(http.MethodGet) // 指定された問題IDの解答を取得

	// ユーザーに対する解答の取得
	publicRoutes.HandleFunc("/users/{user_id}/solutions", handlers.GetSolutionsByUserIDHandler(db)).Methods(http.MethodGet) // ユーザーIDに基づく解答の取得

	// ルートURLのハンドラーを設定
	publicRoutes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the ProCon Web Service :)")
	})

	// 2. 認証が必要なAPIルート
	authRoutes := router.PathPrefix("/api").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware)

	// 問題に関するAPI
	authRoutes.HandleFunc("/problems", handlers.UploadProblemHandler(db)).Methods(http.MethodPost) // 問題の投稿(認証が必要)
	// 解答に関するAPI
	authRoutes.HandleFunc("/problems/{problem_id}/solutions", handlers.SubmitSolutionHandler(db)).Methods(http.MethodPost) // 解答の提出(認証が必要)
	// ユーザーに関するAPI
	authRoutes.HandleFunc("/users/logout", handlers.LogoutUserHandler(db)).Methods(http.MethodPost) // ログアウト(認証が必要)

	// 3. より詳細な権限設定が必要なAPIルート
	authRoutes.HandleFunc("/problems/{problem_id}", middleware.ProblemOwnerMiddlewareFactory(db)(handlers.UpdateProblemHandler(db))).Methods(http.MethodPut)    // 問題の更新(problem_idが必要 + 問題の所有者のみ)
	authRoutes.HandleFunc("/problems/{problem_id}", middleware.ProblemOwnerMiddlewareFactory(db)(handlers.DeleteProblemHandler(db))).Methods(http.MethodDelete) // 問題の削除(problem_idが必要 + 問題の所有者のみ)
	// ユーザーに関するAPI
	authRoutes.HandleFunc("/users/{user_id}", middleware.UserAuthMiddlewareFactory(db)(handlers.UpdateUserHandler(db))).Methods(http.MethodPut) // ユーザープロファイルの更新(user_idが必要 + ユーザー自身のみ)

	// WebSocket通信用のルーティング
	publicRoutes.HandleFunc("/ws", handlers.WebSocketHandler(db))
}
