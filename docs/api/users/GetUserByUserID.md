# `/api/users/{user_id}` (GET): ユーザーID指定でのユーザープロファイルの取得

## 概要:
指定されたユーザーIDに基づいて，ユーザープロファイル情報を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/users/{user_id}`

## URLパラメータ:
- `user_id`: 取得したいユーザーのID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ:
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 指定したユーザーのプロファイル情報
```json
{
    "message": null,
    "result": {
        "user_id": 1,
        "username": "testuser1"
    },
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "User with UserID 1 not found",
    "result": null,
    "status": 404
}  
```

## テスト用curlコマンドの例

```json
curl -X GET http://localhost:8080/api/users/1

{
    "message": null,
    "result": {
        "user_id": 1,
        "username": "testuser1"
    },
    "status": 200
}
```

```json
curl -X GET http://localhost:8080/api/users/2  

{
    "message": "User with UserID 2 not found",
    "result": null,
    "status": 404
}
```