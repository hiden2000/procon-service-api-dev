# `/api/solutions/{solution_id}/result` (GET): 解答判定結果の取得

## 概要:
指定された解答IDに基づいて，提出された解答のサーバー上での判定結果を取得する．

## HTTPメソッド:
GET

## URL構造:
`/api/solutions/{solution_id}/result`

## URLパラメータ:
- `solution_id`: 指定された解答ID

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
不要

## リクエストボディ:
不要

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 解答判定結果の詳細
```json
{
    "message": null,
    "result": {
        "total_cases": 4,
        "correct_cases": 4,
        "incorrect_cases": 0,
        "time_limit_exceeded": 0,
        "case_results": [
            {
                "case_name": "case01.txt",
                "result": "PASSED",
                "execution_time": 62187087
            },
            {
                "case_name": "case02.txt",
                "result": "PASSED",
                "execution_time": 93126322
            },
            {
                "case_name": "case03.txt",
                "result": "PASSED",
                "execution_time": 125165652
            },
            {
                "case_name": "case04.txt",
                "result": "PASSED",
                "execution_time": 93810441
            }
        ]
    },
    "status": 200
}
```
## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "ResultDetails with SolutionID 2 not found",
    "result": null,
    "status": 404
}
```

エラーメッセージ（例）
```json
{
    "message": "variable solution_id error: strconv.Atoi: parsing \"invalid\": invalid syntax",
    "result": null,
    "status": 400
}
```

## テスト用curlコマンドの例

```json
curl -X GET http://localhost:8080/api/solutions/1/result

{
    "message": null,
    "result": {
        "total_cases": 4,
        "correct_cases": 4,
        "incorrect_cases": 0,
        "time_limit_exceeded": 0,
        "case_results": [
            {
                "case_name": "case01.txt",
                "result": "PASSED",
                "execution_time": 62187087
            },
            {
                "case_name": "case02.txt",
                "result": "PASSED",
                "execution_time": 93126322
            },
            {
                "case_name": "case03.txt",
                "result": "PASSED",
                "execution_time": 125165652
            },
            {
                "case_name": "case04.txt",
                "result": "PASSED",
                "execution_time": 93810441
            }
        ]
    },
    "status": 200
}
```