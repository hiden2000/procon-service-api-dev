# `api/users/{user_id}/solutions` (GET): 特定ユーザーの解答を取得

## 概要:
指定されたユーザIDに基づいて，特定のユーザーによって提出された解答を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/users/{user_id}/solutions`

## URLパラメータ:
- `user_id`: 解答を取得したいユーザーのID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ:
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: ユーザーによる解答の一覧
```json
{
    "message": null,
    "result": [
        {
            "solution_id": 1,
            "user_id": 1,
            "problem_id": 1,
            "language_id": 1,
            "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
            "submitted_at": "2024-02-25T07:54:31Z"
        }
    ],
    "status": 200
}
```

## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "variable user_id error: strconv.Atoi: parsing \"invalid\": invalid syntax",
    "result": null,
    "status": 400
}
```


## テスト用curlコマンドの例

```json
curl -X GET "http://localhost:8080/api/users/1/solutions"

{
    "message": null,
    "result": [
        {
            "solution_id": 1,
            "user_id": 1,
            "problem_id": 1,
            "language_id": 1,
            "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
            "submitted_at": "2024-02-25T07:54:31Z"
        }
    ],
    "status": 200
}
```