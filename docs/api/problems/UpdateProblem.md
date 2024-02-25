# `/api/problem/{problem_id}` (PUT): 問題の更新

## 概要:
特定の問題IDに基づいて，問題の内容を更新

## HTTPメソッド:
PUT

## URL構造:
`/api/problem/{problem_id}`

## URLパラメータ:
- `problem_id`: 更新したい問題のID

## 認証用リクエストヘッダー
必要

## リクエストボディ:
- マルチパートフォームデータ
- `title`: 問題のタイトル（必須）
- `description`: 問題の説明（任意）
- `difficulty`: 難易度（必須）
- `input_file`: アップロードする入力ファイル（任意）
- `output_file`: アップロードする出力ファイル（任意）
- 制約として，input_fileに対応する入力ファイル名とoutput_fileに対応する出力ファイルのファイル名は一対一に対応しなくてはいけない
- また，それぞれ重複した名前は許さない．

```json
{
  "title": "Updated Example Problem",
  "description": "Updated description of problem.",
  "input_format": "Updated input description",
  "output_format": "Updated output description",
  "sample_io": "Updated example input and output",
  "difficulty": 2
}
```

## 成功時のレスポンス:
- HTTPステータスコード: 200 OK
- レスポンスボディ: 更新された問題の詳細情報

TODO: データベースに存在しない問題に対して，Update操作をかけたら何が起きるかを書く

```json
{
  "problem_id": 123,
  "title": "Updated Example Problem",
  "description": "Updated description of problem.",
  "input_format": "Updated input description",
  "output_format": "Updated output description",
  "sample_io": "Updated example input and output",
  "difficulty": 2,
  "created_at": "2021-01-01T00:00:00Z",
  "updated_at": "2021-01-01T00:00:00Z"
}
```

## エラー時のレスポンス:
- HTTPステータスコード: 400 Bad Request

エラーメッセージ（例）
```json
{
  "error": "Invalid input data"
}
```

## テスト用curlコマンドの例 

```json
curl -X PUT http://localhost:8080/api/problem/1 \
-H "Content-Type: application/json" \
-d '{
  "title": "updatedproblemtitle20240132",
  "description": "Updated description here.",
  "input_format": "Updated input format.",
  "output_format": "Updated output format.",
  "sample_io": "Updated sample IO.",
  "difficulty": 5,
}'

# TODO: 返却値のidの値が，初期値のまま返ってきてる
>> {"problem_id":1,"title":"updatedproblemtitle","description":"Updated description here.","input_format":"","output_format":"","sample_io":"","difficulty":5,"category":"Updated category","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","io_files":null}
```