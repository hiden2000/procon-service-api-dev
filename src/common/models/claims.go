package models

import "github.com/dgrijalva/jwt-go"

// Claimsは，JWT認証に使用されるクレーム情報を保持する構造体である．
type Claims struct {
	Username string `json:"username"` // ユーザー名である．
	UserID   int    `json:"user_id"`  // ユーザーの一意識別子である．
	jwt.StandardClaims
}
