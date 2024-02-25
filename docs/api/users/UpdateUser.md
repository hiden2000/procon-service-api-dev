# `/api/users/{user_id}` (PUT): ユーザープロファイルの更新

## 概要:
指定されたユーザーIDに基づいて，ユーザープロファイル情報を更新する．

## HTTPメソッド:
PUT

## URL構造:
`/api/users/{user_id}`

## URLパラメータ:
- `user_id`: 更新したいユーザーのID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
必要(ユーザー本人のJWTトークンのみ有効)

## リクエストボディ:
- `username`: 新しいユーザーネーム（必須）

```json
{
  "username": "newusername",
}
```

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 更新されたユーザープロファイル情報
```json
{
    "message": null,
    "result": {
        "user_id": 1,
        "username": "testuser0"
    },
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "transaction execute error: Error 1406 (22001): Data too long for column 'Username' at row 1",
    "result": null,
    "status": 500
}
```

## テスト用curlコマンドの例

```json
curl -X PUT http://localhost:8080/api/users/1 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODUxNzM2fQ.I7QYEUT8AYhq4jsqxLpD3AUfsAOVfHUNJNRcUv2th7Y" \
-d '{"username": "testuser0"}'

{
    "message": null,
    "result": {
        "user_id": 1,
        "username": "testuser0"
    },
    "status": 200
}
```