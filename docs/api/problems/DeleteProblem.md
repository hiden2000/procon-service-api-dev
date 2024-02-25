# `/api/problem/{problem_id}` (DELETE): 問題の削除

## 概要:
特定の問題IDを持つ問題を削除する．

## HTTPメソッド:
DELETE

## URL構造:
`/api/problem/{problem_id}`

## URLパラメータ:
- `problem_id`: 削除する問題のID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
必要

## リクエストボディ:
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: なし
```json

```

## エラー時のレスポンス:

不正なリクエストの場合のエラーメッセージ(例)
```json
{
    "message": "TokenValidationFailed: Token validation failed",
    "result": null,
    "status": 401
}
```

## テスト用curlコマンドの例

```json

curl -X POST http://localhost:8080/api/problems \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODUxMjAyfQ.fsX_y6s2R-fnENGONqyyNdnnY267YeJx50s3EUaQGD4" \
  -F "metadata={ \
        \"title\": \"this is simple a + b problem (3)\", \
        \"description\": \"This is a test problem description. (3)\", \
        \"difficulty\": 1 \
      }" \
  -F "input_file=@/Users/example_user/example_problem/problems/problem1/in/case01.txt" \
  -F "input_file=@/Users/example_user/example_problem/problems/problem1/in/case02.txt" \
  -F "input_file=@/Users/example_user/example_problem/problems/problem1/in/case03.txt" \
  -F "input_file=@/Users/example_user/example_problem/problems/problem1/in/case04.txt" \
  -F "output_file=@/Users/example_user/example_problem/problems/problem1/out/case01.txt" \
  -F "output_file=@/Users/example_user/example_problem/problems/problem1/out/case02.txt" \
  -F "output_file=@/Users/example_user/example_problem/problems/problem1/out/case03.txt" \
  -F "output_file=@/Users/example_user/example_problem/problems/problem1/out/case04.txt"

{
    "message": null,
    "result": {
        "problem_id": 3,
        "user_id": 1,
        "title": "this is simple a + b problem (3)",
        "description": "This is a test problem description. (3)",
        "difficulty": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z",
        "category_ids": null
    },
    "status": 201
}

curl -X DELETE http://localhost:8080/api/problems/4 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODUxMjAyfQ.fsX_y6s2R-fnENGONqyyNdnnY267YeJx50s3EUaQGD4"

```