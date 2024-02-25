# `/api/problems` (GET): 問題の取得

## 概要
このエンドポイントは登録されているすべての問題を取得するために使用される．

## HTTPメソッド
GET

## URL構造
`/api/problems`

## URLパラメータ:
不要

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ
不要

## 成功時のレスポンス
- HTTPステータスコード: 200 OK

レスポンスボディ: 取得された問題のリスト
```json

[
  {
    "problem_id":1,
    "user_id":1,
    "title":"this is first uploaded problem",
    "description":"This is an example problem description.","difficulty":1,"created_at":"2024-02-09T15:08:09Z","updated_at":"2024-02-09T15:08:09Z",
    "category_ids":null
  },
  {
    "problem_id":2,
    "user_id":3,
    "title":"this problem is uploaded by testuser ",
    "description":"This is an example problem description.","difficulty":5,"created_at":"2024-02-09T15:43:11Z","updated_at":"2024-02-09T15:43:11Z",
    "category_ids":null
}
]
```

## エラー時のレスポンス
このエンドポイントでは，特定のエラー条件に基づくレスポンスは予定されていない．

しかし，サーバーやデータベースの問題により，500 Internal Server Error が発生する可能性がある．

## テスト用curlコマンドの例 

```json
curl -X GET http://localhost:8080/api/problems

{
    "message": null,
    "result": [
        {
            "problem_id": 1,
            "user_id": 1,
            "title": "this is simple a + b problem",
            "description": "This is a test problem description.",
            "difficulty": 1,
            "created_at": "2024-02-25T07:32:33Z",
            "updated_at": "2024-02-25T07:32:33Z"
        },
        {
            "problem_id": 2,
            "user_id": 1,
            "title": "this is simple a + b problem (2) ",
            "description": "This is a test problem description (2).",
            "difficulty": 1,
            "created_at": "2024-02-25T07:33:49Z",
            "updated_at": "2024-02-25T07:33:49Z"
        }
    ],
    "status": 200
}
```