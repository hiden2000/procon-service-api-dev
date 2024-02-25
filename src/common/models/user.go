package models

import "time"

// Userは，サービスを利用するユーザーの基本情報を保持する構造体である．
type User struct {
	UserID    int       `json:"user_id"`    // ユーザーの一意識別子である．
	Username  string    `json:"username"`   // ユーザー名である．
	Password  string    `json:"password"`   // ユーザーのパスワード（ハッシュ化される）である．
	CreatedAt time.Time `json:"created_at"` // ユーザーのアカウント作成日時である．
	LastLogin time.Time `json:"last_login"` // ユーザーの最終ログイン日時である．
}

// UserCredentialsは，ユーザー認証時に使用される認証情報を保持する構造体である．
// ユーザー名とパスワードが含まれる．
type UserCredentials struct {
	Username string `json:"username"` // ユーザー名である．
	Password string `json:"password"` // パスワードである．
}

// UserProfileは，ユーザーのプロファイル情報を表す構造体である．
// ユーザーIDとユーザー名が含まれる．
type UserProfile struct {
	UserID   int    `json:"user_id"`  // ユーザーの一意識別子である．
	Username string `json:"username"` // ユーザー名である．
}
