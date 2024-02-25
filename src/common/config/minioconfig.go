package config

import "os"

// MinIOConfigはMinIOオブジェクトストレージに接続するための設定を保持する構造体である．
// 環境変数から設定値を読み込み，アプリケーションからMinIOサーバーにアクセスする際に使用される．
type MinIOConfig struct {
	Endpoint   string // MinIOサーバーのエンドポイント．
	User       string // MinIOサーバーへのアクセスに使用するユーザー名．
	Password   string // MinIOサーバーへのアクセスに使用するパスワード．
	UseSSL     bool   // SSLを使用してMinIOサーバーに接続するかどうか．
	BucketName string // 使用するバケットの名前．
}

// NewMinIOConfigはMinIOConfigの新しいインスタンスを生成し，環境変数から設定値を読み込んで返す関数である．
func NewMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:   os.Getenv("MINIO_ENDPOINT"), //"minio:9000"
		User:       os.Getenv("MINIO_ROOT_USER"),
		Password:   os.Getenv("MINIO_ROOT_PASSWORD"),
		UseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
		BucketName: os.Getenv("MINIO_BUCKET_NAME"),
	}
}
