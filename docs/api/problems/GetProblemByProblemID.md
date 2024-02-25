# `/api/problems/{problem_id}` (GET): 特定の問題の取得

## 概要:
特定の問題IDに基づいて，指定された問題の詳細情報を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/problems/{problem_id}`

## URLパラメータ:
- `problem_id`: 取得したい問題のID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 指定した問題の詳細情報
```json
{
    "message": null,
    "result": {
        "problem_id": 1,
        "user_id": 1,
        "title": "this is simple a + b problem",
        "description": "This is a test problem description.",
        "difficulty": 1,
        "created_at": "2024-02-25T07:32:33Z",
        "updated_at": "2024-02-25T07:32:33Z",
        "category_ids": null
    },
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ (例)
```json
{
    "message": "Problem with ProblemID 10290 not found",
    "result": null,
    "status": 404
}
```

## テスト用curlコマンドの例

```json
curl -X GET http://localhost:8080/api/problems/1

{
    "message": null,
    "result": {
        "problem_id": 1,
        "user_id": 1,
        "title": "this is simple a + b problem",
        "description": "This is a test problem description.",
        "difficulty": 1,
        "created_at": "2024-02-25T07:32:33Z",
        "updated_at": "2024-02-25T07:32:33Z",
        "category_ids": null
    },
    "status": 200
}
```