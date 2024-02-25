# `/api/users/login` (POST): ユーザーログイン

## 概要:
ユーザーがログインし，JWTトークンとユーザーの基本情報を取得する．

生成されたJWTトークンはRedisに保存され，今後のセッション管理に利用される．

## HTTPメソッド:
POST

## URL構造:
`/api/users/login`

## URLパラメータ:
不要

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ:
ユーザーのログイン情報を含むJSONオブジェクト

- `username`: ユーザー名（必須）
- `password`: パスワード（必須）

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: ログインに成功したユーザーのJWTトークンとユーザー情報
```json
{
    "message": "",
    "result": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA3ODQ1NzQ5fQ.P6OUiyZezV1aayV7gVfYVax54Jz2qt_qP7tjvSKKMIc",
        "user": {
            "user_id": 1,
            "username": "testuser",
        }
    },
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "パスワードが不正です．",
    "result": null,
    "status": 401
}
```

エラーメッセージ（例）
```json
{
    "message": "User with Username testuser1 not found",
    "result": null,
    "status": 404
}
```

## テスト用curlコマンドの例
```json
curl -X POST http://localhost:8080/api/users/login \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "password123"
         }'
{
    "message": "パスワードが不正です．",
    "result": null,
    "status": 401
}
```
```json
curl -X POST http://localhost:8080/api/users/login \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser1",
           "password": "password123"
         }'
{
    "message": "User with Username testuser1 not found",
    "result": null,
    "status": 404
}
```
```json
curl -X POST http://localhost:8080/api/users/login \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "testpassword"
         }'
{
    "message": "",
    "result": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA3ODQ1NzQ5fQ.P6OUiyZezV1aayV7gVfYVax54Jz2qt_qP7tjvSKKMIc",
        "user": {
            "user_id": 1,
            "username": "testuser",
        }
    },
    "status": 200
}
```