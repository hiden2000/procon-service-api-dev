package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"procon_web_service/src/common/middleware"
	"procon_web_service/src/web/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

var db *sql.DB // データベース接続を保持するグローバル変数

func initDB() *sql.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbname)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func main() {
	db = initDB()
	defer db.Close()

	// マルチプレクサーの作成
	router := mux.NewRouter()

	// APIルーティングの設定 && データベース接続
	routes.RegisterApiRoutes(router, db)

	// CORSの設定
	handler := cors.AllowAll().Handler(router)

	limiter := rate.NewLimiter(5, 5) // 1秒あたり5リクエストまで許可し，バーストサイズも5に設定

	// レート制限ミドルウェアの適用
	limitedHandler := middleware.RateLimiterMiddleware(limiter)(handler)

	// HTTPサーバーの開始
	if err := http.ListenAndServe(":8080", limitedHandler); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
