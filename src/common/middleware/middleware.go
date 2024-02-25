package middleware

import (
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

// LoggingMiddlewareはリクエストのメソッドとURLをログに記録するミドルウェアである．
// HTTPリクエストを受け取ると，そのメソッドとURLをログ出力し，次のミドルウェアまたはハンドラへリクエストを渡す．
// ログは開発やデバッグ時にリクエストの流れを追跡しやすくするために役立つ．
//
// パラメータ:
// - next http.Handler: リクエスト処理を続けるための次のハンドラまたはミドルウェア．
//
// 戻り値:
// - http.Handler: ログ記録機能を備えたHTTPリクエストのハンドラ．
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// RateLimiterMiddlewareはリクエストのレート制限を実施するミドルウェアである．
// このミドルウェアは，指定されたレートリミッターを使用してリクエストの処理速度を制限し，
// システムへの過負荷を防ぐ．レートリミットを超えたリクエストは，503 Service Unavailableエラーとして拒否される．
//
// パラメータ:
// - limiter *rate.Limiter: リクエストのレートを制限するためのレートリミッター．
//
// 戻り値:
// - func(http.Handler) http.Handler: レート制限機能を備えたミドルウェアを生成する関数．
func RateLimiterMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
