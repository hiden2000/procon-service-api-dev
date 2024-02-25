# `api/problems/{problem_id}/solutions` (POST): 解答の提出

## 概要:
このエンドポイントはユーザーが特定の問題に対する解答を提出するために使用される．

## HTTPメソッド:
POST

## URL構造:
`api/problems/{problem_id}/solutions`

## URLパラメータ:
- `problem_id`: 解答を提出したい問題のID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
必要

## リクエストボディ:
- `language_id`: 解答の言語のID（必須）
- `code`: 解答コード（必須）

```json
{
  "language_id": 1,
  "code": "print('Hello, World!')"
}
```

## 成功時のレスポンス:
- HTTPステータスコード: 201 Created

レスポンスボディ: 提出された解答の詳細
```json
{
    "message": null,
    "result": {
        "solution_id": 1,
        "user_id": 1,
        "problem_id": 1,
        "language_id": 1,
        "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
        "submitted_at": "0001-01-01T00:00:00Z"
    },
    "status": 201
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "TokenValidationFailed: Token validation failed",
    "result": null,
    "status": 401
}
```

## テスト用curlコマンドの例

```json
curl -X POST http://localhost:8080/api/problems/1/solutions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODU3MDgwfQ.gCYDC8rkXnuy3OULx9GFxSi77qjvRDKSlz-V_QMkkcg" \
-d '{
  "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
  "language_id": 1
}'

{
    "message": null,
    "result": {
        "solution_id": 1,
        "user_id": 1,
        "problem_id": 1,
        "language_id": 1,
        "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
        "submitted_at": "0001-01-01T00:00:00Z"
    },
    "status": 201
}
```

