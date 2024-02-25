package config

import (
	"log"
	"os"
	"time"
)

// JWTConfigは，JWTトークンの生成と検証に関連する設定を保持する構造体である．
// これには，シークレットキーとトークンの有効期限を含む設定項目が含まれる．
//
// フィールド:
// - SecretKey string: JWTトークンの署名に使用するシークレットキー．
// - ExpirationTime time.Duration: トークンの有効期限．
type JWTConfig struct {
	SecretKey      string
	ExpirationTime time.Duration
}

// NewJWTConfigは，環境変数からJWT設定を読み込み，JWTConfigインスタンスを生成する関数である．
// 戻り値として，初期化されたJWTConfigのポインタを返す．
func NewJWTConfig() *JWTConfig {
	// 環境変数一覧をログに出力
	for _, env := range os.Environ() {
		log.Println("デバッグ", env)
	}
	return &JWTConfig{
		SecretKey:      os.Getenv("JWT_SECRET_KEY"),
		ExpirationTime: 2 * time.Hour, // 例として2時間をデフォルト値として設定
	}
}
