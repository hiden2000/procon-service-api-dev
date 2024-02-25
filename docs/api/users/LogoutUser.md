# `/api/user/logout` (POST): ユーザーログアウト

## 概要:
ユーザーがサービスからログアウトする．

## HTTPメソッド:
POST

## URL構造:
`/api/user/logout`

## URLパラメータ:
不要

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
必要

## リクエストボディ:
なし

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: ログアウト成功のメッセージ
```json
{
    "message": null,
    "result": null,
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "TokenExpired: Token expired or not active yet",
    "result": null,
    "status": 401
}

{
    "message": "InvalidTokenFormat: Invalid token format",
    "result": null,
    "status": 401
}
```

## テスト用curlコマンドの例

```json
curl -X POST http://localhost:8080/api/users/logout \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA3ODIzODgzfQ.mts3skdDPmjbAiP0PDnc7cC1zSDY3sirB9MAB16h938"

{
    "message": "TokenExpired: Token expired or not active yet",
    "result": null,
    "status": 401
}
```

```json
curl -X POST http://localhost:8080/api/users/logout \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA3ODIzODgzfDPmjbAiP0PDnc7cC1zSDY3sirB9MAB16h9" 

{
    "message": "InvalidTokenFormat: Invalid token format",
    "result": null,
    "status": 401
}
```