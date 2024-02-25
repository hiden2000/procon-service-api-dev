package utils

import (
	"context"
	"log"
	"net/http"
	commonerrors "procon_web_service/src/common/errors"
	"procon_web_service/src/common/models"
	"procon_web_service/src/web/config"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ctx = context.Background() // Redis操作のためのコンテキスト

	jwtKey []byte // JWTシークレットキー

	durationTime = 2 * time.Hour // セッション期限
)

func init() {
	jwtconfig := config.NewJWTConfig()
	if jwtconfig.SecretKey == "" {
		log.Fatal("${JWT_SECRET_KEY} is not defined")
	}
	// JWTシークレットキーとセッション期限を設定
	jwtKey = []byte(jwtconfig.SecretKey)
	durationTime = jwtconfig.ExpirationTime
}

// GenerateJWTは，指定されたユーザー情報を基にJWTトークンを生成し，Redisサーバーへ保存してセッションを有効化する．
// 成功時には生成されたJWTトークン文字列を返し，エラーが発生した場合にはエラーを返す．
//
// パラメータ:
// - user models.User: トークン生成の基となるユーザー情報．
//
// 戻り値:
// - string: 生成されたJWTトークン文字列．
// - error: トークン生成やRedisへの保存に失敗した場合のエラー．
func GenerateJWT(user models.User) (string, error) {
	expirationTime := time.Now().Add(durationTime)
	claims := &models.Claims{
		Username: user.Username,
		UserID:   user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", commonerrors.WrapTokenError(err)
	}

	// Redisにトークンを保存してセッションを有効化
	err = redisClient.Set(ctx, tokenString, user.UserID, time.Until(expirationTime)).Err()
	if err != nil {
		return "", commonerrors.WrapTokenError(err)
	}

	return tokenString, nil
}

// VerifyPasswordは，提供されたパスワードがハッシュ化されたパスワードと一致するか検証する．
// 検証に成功した場合はnilを，失敗した場合はエラーを返す．
//
// パラメータ:
// - hash string: 検証対象のハッシュ化されたパスワード．
// - password string: 検証する生パスワード．
//
// 戻り値:
// - error: パスワードが一致しない場合のエラー．
func VerifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return commonerrors.NewAccessDeniedError("パスワードが不正です．")
	}
	return nil
}

// GenerateHashedPasswordは，提供されたパスワードをbcryptアルゴリズムを使用してハッシュ化する．
// 成功時にはハッシュ化されたパスワードを文字列で返し，エラーが発生した場合にはエラーを返す．
//
// パラメータ:
// - password string: ハッシュ化する生パスワード．
//
// 戻り値:
// - string: ハッシュ化されたパスワード．
// - error: ハッシュ化処理に失敗した場合のエラー．
func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", commonerrors.WrapTokenError(err)
	}
	return string(hashedPassword), nil
}

// IsUserAuthenticatedは，HTTPリクエストから抽出したJWTトークンを検証し，ユーザーが認証されているか確認する．
// 認証が成功した場合はユーザーのクレーム情報を返し，失敗した場合はエラーを返す．
//
// パラメータ:
// - r *http.Request: 認証を検証するHTTPリクエスト．
//
// 戻り値:
// - *models.Claims: 認証に成功したユーザーのクレーム情報．
// - error: 認証に失敗した場合のエラー．
func IsUserAuthenticated(r *http.Request) (*models.Claims, error) {
	tokenString := ExtractToken(r)

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, commonerrors.WrapTokenError(err)
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		isValid, err := isSessionValid(claims.UserID, tokenString)
		if err != nil {
			return nil, commonerrors.WrapTokenError(err)
		}
		if isValid {
			return claims, nil
		}
		return nil, commonerrors.NewTokenError("SessionInvalid", "セッションが無効です．")
	}

	return nil, commonerrors.NewTokenError("TokenInvalid", "トークンが無効です．")
}

// IsTokenAuthenticatedは，渡されたJWTトークンを検証し，ユーザーが認証されているか確認する．
// 認証が成功した場合はユーザーのクレーム情報を返し，失敗した場合はエラーを返す．
//
// パラメータ:
// - tokenString string: 認証を検証するトークン．．
//
// 戻り値:
// - *models.Claims: 認証に成功したユーザーのクレーム情報．
// - error: 認証に失敗した場合のエラー．
func IsTokenAuthenticated(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, commonerrors.WrapTokenError(err)
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		isValid, err := isSessionValid(claims.UserID, tokenString)
		if err != nil {
			return nil, commonerrors.WrapTokenError(err)
		}
		if isValid {
			return claims, nil
		}
		return nil, commonerrors.NewTokenError("SessionInvalid", "セッションが無効です．")
	}

	return nil, commonerrors.NewTokenError("TokenInvalid", "トークンが無効です．")
}

// ExtractTokenは，HTTPリクエストのAuthorizationヘッダからJWTトークンを抽出する．
// トークンが存在する場合はその文字列を，存在しない場合は空文字を返す．
//
// パラメータ:
// - r *http.Request: トークンを抽出するHTTPリクエスト．
//
// 戻り値:
// - string: 抽出されたJWTトークン文字列．
func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// EraseTokenは，指定されたJWTトークンをRedisから削除し，セッションを無効化する．
// 処理に成功した場合はnilを，失敗した場合はエラーを返す．
//
// パラメータ:
// - tokenString string: 削除するJWTトークン文字列．
//
// 戻り値:
// - error: トークン削除処理に失敗した場合のエラー．
func EraseToken(tokenString string) error {
	err := redisClient.Del(ctx, tokenString).Err()
	if err != nil {
		return commonerrors.WrapTokenError(err)
	}
	return nil
}

// isSessionValid - セッションが有効かどうかを確認
func isSessionValid(userID int, tokenString string) (bool, error) {
	result, err := redisClient.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		return false, commonerrors.NewTokenError("SessionNotFound", "セッションが見つかりません．")
	} else if err != nil {
		// Redisのエラーをラップして返す
		return false, commonerrors.WrapTokenError(err)
	}
	return result == strconv.Itoa(userID), nil
}
