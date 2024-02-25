# `/api/solutions/{solution_id}` (GET): 特定の解答の詳細情報の取得

## 概要:
指定された解答IDに基づいて，特定の解答の詳細情報を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/solutions/{solution_id}`

## URLパラメータ:
- `solution_id`: 取得したい解答のID

## クエリパラメータ
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ:
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 特定の解答の詳細情報
```json
{
    "message": null,
    "result": {
        "solution_id": 1,
        "user_id": 1,
        "problem_id": 1,
        "language_id": 1,
        "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
        "submitted_at": "2024-02-25T07:54:31Z"
    },
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "Solution with SolutionID 12090 not found",
    "result": null,
    "status": 404
}
```

## テスト用curlコマンドの例

```json
curl -X GET http://localhost:8080/api/solutions/1   
 
{
    "message": null,
    "result": {
        "solution_id": 1,
        "user_id": 1,
        "problem_id": 1,
        "language_id": 1,
        "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
        "submitted_at": "2024-02-25T07:54:31Z"
    },
    "status": 200
}
```