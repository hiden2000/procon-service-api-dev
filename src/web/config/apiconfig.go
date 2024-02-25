package config

import "os"

// APIEndpointsConfigは，外部APIエンドポイントの設定を保持する構造体である．
// これには，ジャッジサーバーのURLを含む設定項目が含まれる．
//
// フィールド:
// - JudgeServerURL string: ジャッジサーバーのベースURL．
type APIEndpointsConfig struct {
	JudgeServerURL string
}

// NewAPIEndpointsConfigは，環境変数からAPIエンドポイントの設定を読み込み，APIEndpointsConfigインスタンスを生成する関数である．
// 戻り値として，初期化されたAPIEndpointsConfigのポインタを返す．
func NewAPIEndpointsConfig() *APIEndpointsConfig {
	return &APIEndpointsConfig{
		JudgeServerURL: os.Getenv("JUDGE_SERVER_URL"),
	}
}
