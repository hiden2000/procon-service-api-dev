package config

import "os"

// DatabaseConfigは，データベース接続の設定を保持する構造体である．
// これには，ホスト，ユーザー名，パスワード，データベース名を含む設定項目が含まれる．
//
// フィールド:
// - Host string: データベースサーバーのホスト名またはIPアドレス．
// - User string: データベース接続に使用するユーザー名．
// - Password string: データベース接続に使用するパスワード．
// - Database string: 接続するデータベースの名前．
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Database string
}

// NewDatabaseConfigは，環境変数からデータベース接続の設定を読み込み，DatabaseConfigインスタンスを生成する関数である．
// 戻り値として，初期化されたDatabaseConfigのポインタを返す．
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}
}
