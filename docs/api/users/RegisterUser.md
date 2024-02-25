# `/api/users` (POST): ユーザー登録

## 概要:
新規ユーザーを登録し，今後のログインに必要なJWTトークンを返却する．

生成されたJWTトークンはサービス内部のRedisに保存され，以降のセッション間に使用される．

## HTTPメソッド:
POST

## URL構造:
`/api/users`

## URLパラメータ:
不要

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ:
新規ユーザーの情報を含むJSONオブジェクト

- `username`: ユーザー名（必須）
- `password`: パスワード（必須）

## 成功時のレスポンス:
- HTTPステータスコード: 201 Created

レスポンスボディ: 生成されたJWTトークン等を含めた，ユーザーの登録情報
```json
{
    "message": null,
    "result": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODUxMjAyfQ.fsX_y6s2R-fnENGONqyyNdnnY267YeJx50s3EUaQGD4",
        "user": {
            "user_id": 1,
            "username": "testuser",
            "password": "$2a$10$2r.ETu2Tpxv0Vh.HwsWJCOAO50sUxos6Z5dQo1QZnwlvc.kmj5Dia",
            "created_at": "0001-01-01T00:00:00Z",
            "last_login": "0001-01-01T00:00:00Z"
        }
    },
    "status": 201
}

```
## エラー時のレスポンス:

ユーザー名，メールアドレス，またはパスワードなどが無効な場合のエラーメッセージ
```json
{
    "message": "constraint violation: username, This username is already in use.",
    "result": null,
    "status": 400
}

{
    "message": "Failed to decode request body: EOF",
    "result": null,
    "status": 400
}
```

## テスト用curlコマンドの例

```json
curl -X POST http://localhost:8080/api/users \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "testpassword"
         }'

{
    "message": null,
    "result": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODUxMjAyfQ.fsX_y6s2R-fnENGONqyyNdnnY267YeJx50s3EUaQGD4",
        "user": {
            "user_id": 1,
            "username": "testuser",
            "password": "$2a$10$2r.ETu2Tpxv0Vh.HwsWJCOAO50sUxos6Z5dQo1QZnwlvc.kmj5Dia",
            "created_at": "0001-01-01T00:00:00Z",
            "last_login": "0001-01-01T00:00:00Z"
        }
    },
    "status": 201
}
```