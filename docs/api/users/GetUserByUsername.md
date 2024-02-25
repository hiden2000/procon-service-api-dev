# `/api/users` (GET): ユーザー名指定でのユーザープロファイルの取得

## 概要:
指定されたユーザー名に基づいて，ユーザープロファイル情報を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/users`

## URLパラメータ:
不要

## クエリパラメータ:
- `username`: 取得したいユーザーの名前

## 認証用リクエストヘッダー
不要

## リクエストボディ:
なし

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンス内容: 指定したユーザーのプロファイル情報
```json
{
    "message": "",
    "result": {
        "user_id": 1,
        "username": "testuser",
    },
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "User with Username notfound not found",
    "result": null,
    "status": 404
} 
```

## テスト用curlコマンドの例

```json
curl -X GET "http://localhost:8080/api/users?username=testuser"
{
    "message": null,
    "result": {
        "user_id": 4,
        "username": "testuser1"
    },
    "status": 200
}
```

```json
curl -X GET "http://localhost:8080/api/users?username=notfound"
{
    "message": "User with Username notfound not found",
    "result": null,
    "status": 404
}                                                                           
```