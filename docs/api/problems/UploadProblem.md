# `/api/problems` (POST): 問題の投稿（ファイルアップロード含む）

## 概要:
このエンドポイントは新しい問題のデータと関連するファイルをアップロードし，問題を作成する．

## HTTPメソッド:
POST

## URL構造:
`/api/problems`

## URLパラメータ:
不要

## クエリパラメータ:
不要

## 認証用リクエストヘッダー
必要

## リクエストボディ:
- マルチパートフォームデータ
- `title`: 問題のタイトル（必須）
- `description`: 問題の説明（任意）
- `difficulty`: 難易度（必須）
- `input_file`: アップロードする入力ファイル（任意）
- `output_file`: アップロードする出力ファイル（任意）
- 制約として，input_fileに対応する入力ファイル名とoutput_fileに対応する出力ファイルのファイル名は一対一に対応しなくてはいけない．
- また，それぞれ重複した名前は許さない．

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK

レスポンスボディ: 作成された問題に関する情報
```json
{
    "message": null,
    "result": {
        "problem_id": 1,
        "user_id": 1,
        "title": "this is simple a + b problem",
        "description": "This is a test problem description.",
        "difficulty": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    },
    "status": 201
}
```
## エラー時のレスポンス:

エラーメッセージ（例）
```json
{
    "message": "File validation error: input_file: case04.txt に対応する output_file が存在しません",
    "result": null,
    "status": 400
}
```
エラーメッセージ（例）
```json
{
    "message": "File validation error: input_file: case04.txt に対応する input_file が存在しません",
    "result": null,
    "status": 400
}

```

## テスト用curlコマンドの例 

```json
curl -X POST http://localhost:8080/api/problems \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODU3MDgwfQ.gCYDC8rkXnuy3OULx9GFxSi77qjvRDKSlz-V_QMkkcg" \
  -F "metadata={ \
        \"title\": \"this is simple a + b problem\", \
        \"description\": \"This is a test problem description.\", \
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
        "problem_id": 1,
        "user_id": 1,
        "title": "this is simple a + b problem",
        "description": "This is a test problem description.",
        "difficulty": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z",
    },
    "status": 201
}
```
