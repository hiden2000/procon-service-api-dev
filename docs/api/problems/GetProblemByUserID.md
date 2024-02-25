# `/api/users/{user_id}/problems` (GET): 特定の問題の取得

## 概要:
特定のユーザーIDに基づいて，そのユーザーが作成した問題の詳細情報を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/users/{user_id}/problems`

## URLパラメータ:
- `user_id`: 取得したい問題のID

## 認証用リクエストヘッダー
不要

## リクエストボディ
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 指定された問題群の詳細情報
```json
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
            "updated_at": "2024-02-25T07:32:33Z",
            "category_ids": null
        },
        {
            "problem_id": 2,
            "user_id": 1,
            "title": "this is simple a + b problem (2) ",
            "description": "This is a test problem description (2).",
            "difficulty": 1,
            "created_at": "2024-02-25T07:33:49Z",
            "updated_at": "2024-02-25T07:33:49Z",
            "category_ids": null
        }
    ],
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ (例)
```json
{
    "message": "variable user_id error: strconv.Atoi: parsing \"invalid\": invalid syntax",
    "result": null,
    "status": 400
}
```

## テスト用curlコマンドの例

```json
curl -X GET http://localhost:8080/api/users/1/problems

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
            "updated_at": "2024-02-25T07:32:33Z",
            "category_ids": null
        },
        {
            "problem_id": 2,
            "user_id": 1,
            "title": "this is simple a + b problem (2) ",
            "description": "This is a test problem description (2).",
            "difficulty": 1,
            "created_at": "2024-02-25T07:33:49Z",
            "updated_at": "2024-02-25T07:33:49Z",
            "category_ids": null
        }
    ],
    "status": 200
}
```