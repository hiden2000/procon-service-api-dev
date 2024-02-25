package main

import (
	"log"
	"net/http"
	"procon_web_service/src/common/middleware"
	"procon_web_service/src/judge/routes"
	"procon_web_service/src/judge/utils"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// クリーンアップスケジューラの開始
	utils.StartCleanupScheduler(30*time.Minute, 2*time.Hour) // 30分ごとに実行 && 2時間以上前のファイルを削除

	limiter := rate.NewLimiter(5, 5) // 1秒あたり5リクエストまで許可し、バーストサイズも5に設定
	router := routes.SetupRoutes()

	// レート制限ミドルウェアの適用
	limitedRouter := middleware.RateLimiterMiddleware(limiter)(router)

	if err := http.ListenAndServe(":8080", limitedRouter); err != nil {
		log.Fatal(err)
	}
}
