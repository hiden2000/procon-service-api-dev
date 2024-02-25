# `/api/ws` (WebSocket): 解答判定の非同期通知

## 概要:
このWebSocketエンドポイントは，ユーザーが提出した解答の判定結果を非同期で通知する．

## 通信方式:
WebSocket

## URL構造:
`/api/ws`

## 接続方法:
クライアントはWebSocketプロトコルを使用してこのエンドポイントに接続する．接続後，解答の提出結果の非同期通知を待機する．

なお，提出にあたってwebsocket内部で通信の認証を行うために，URLの内部にJWTトークンを入れることによって，サーバーと認証を行う．

## メッセージ形式:
クライアントからサーバーへのメッセージは，解答のJSON形式のデータを含む．

サーバーからクライアントへのメッセージは，解答の判定結果を含むJSON形式のデータとなる．

クライアントからのメッセージ例:
```json
{
  "solution_id": 1,
  "user_id": 1,
  "problem_id": 1,
  "language_id": 1,
  "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
  "submitted_at": "2024-02-25T07:54:31Z"
}
```

サーバからのメッセージ例:
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
        "case_name": "case04.txt",
        "result": "PASSED",
        "execution_time": 93810441
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
      }
    ]
  },
  "status": 200
}
```
## エラー時の処理:

エラーが発生した場合（例: 判定サーバーへの接続失敗），サーバーはエラーメッセージをクライアントに送信する．

## 使用方法:

- クライアントは /api/solutions エンドポイントに解答をPOSTする．
- 解答の提出後，クライアントは /api/ws に以下に示すようなコマンドで，WebSocket接続を開始する．
- クライアントは解答データをサーバーに送信し，判定結果を待つ．
- サーバーは解答の判定を行い，結果をクライアントに非同期で通知する．
- 実際にサービスとして展開する上では，フロントエンドなどで，非同期通信の取り扱いを行うため，今回のAPIは一般に表には出てこないことに注意する．

## 

# 実際の提出コードの例

```json

websocat "ws://localhost:8080/api/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwidXNlcl9pZCI6MSwiZXhwIjoxNzA4ODU3MDgwfQ.gCYDC8rkXnuy3OULx9GFxSi77qjvRDKSlz-V_QMkkcg"

{
  "solution_id": 1,
  "user_id": 1,
  "problem_id": 1,
  "language_id": 1,
  "code": "X, Y = map(int, input().split())\nprint(X + Y)\n",
  "submitted_at": "2024-02-25T07:54:31Z"
}

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
        "case_name": "case04.txt",
        "result": "PASSED",
        "execution_time": 93810441
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
      }
    ]
  },
  "status": 200
}

```
使用できる言語全てに対して，解答の提出例を以下に示す．(ex. Python, C++, GoLang, Java, Rust)

```json
// Rust
{
  "solution_id": 5,
  "user_id": 10,
  "problem_id": 2,
  "language_id": 5,
  "code": "use std::io::{self, Read};\n\nfn main() {\n    let mut input = String::new();\n    io::stdin().read_to_string(&mut input).unwrap();\n    let nums: Vec<i64> = input.split_whitespace().map(|n| n.parse().unwrap()).collect();\n    let x = nums[0];\n    let y = nums[1];\n    const MOD: i64 = 998244353;\n\n    let mut ans = ((x % MOD) * (y % MOD)) % MOD;\n\n    for i in (x / y)..=x {\n        if i == 0 { continue; }\n        if i * i > x {\n            for j in 1..=y {\n                if x / j >= i {\n                    ans = (ans - ((x / j) % MOD * j) % MOD + MOD) % MOD;\n                } else {\n                    break;\n                }\n            }\n            break;\n        } else {\n            let m = std::cmp::min(x / i, y);\n            let m2 = x / (i + 1);\n            let p_initial = m + m2 + 1;\n            let q_initial = m - m2;\n            let mut p = p_initial;\n            let mut q = q_initial;\n            if p % 2 == 0 {p /= 2;}\n            if q % 2 == 0 {q /= 2;}\n            let r = (p % MOD) * (q % MOD) % MOD;\n            ans = (ans - (r * i % MOD) + MOD) % MOD;\n        }\n    }\n\n    println!(\"{}\", ans);\n}",
  "submitted_at": "2024-02-09T16:26:02Z"
}

// Java
{
  "solution_id": 6,
  "user_id": 10,
  "problem_id": 2,
  "language_id": 4,
  "code": "import java.util.Scanner;\n\nclass Main {\n    public static void main(String[] args) {\n        Scanner sc = new Scanner(System.in);\n        long X = sc.nextLong();\n        long Y = sc.nextLong();\n        final long MOD = 998244353;\n\n        long ans = ((X % MOD) * (Y % MOD)) % MOD;\n\n        for (long i = X / Y; i <= X; i++) {\n            if (i > X / i) {\n                for (long j = 1; j <= Y; j++) {\n                    if (X / j >= i) {\n                        ans -= (X / j) % MOD * (j % MOD) % MOD;\n                        ans = (ans + MOD) % MOD;\n                    } else {\n                        break;\n                    }\n                }\n                break;\n            } else {\n                long M = Math.min(X / i, Y);\n                long m = X / (i + 1);\n                long p = M + m + 1;\n                long q = M - m;\n                if (p % 2 == 0) p /= 2;\n                if (q % 2 == 0) q /= 2;\n                long r = (p % MOD) * (q % MOD) % MOD;\n                r *= i;\n                r %= MOD;\n                ans -= r;\n                ans = (ans + MOD) % MOD;\n            }\n        }\n\n        System.out.println(ans);\n    }\n}",
  "submitted_at": "2024-02-09T16:26:02Z"
}

// GoLang
{
    "solution_id": 7,
    "user_id": 10,
    "problem_id": 2,
    "language_id": 3,
    "code": "package main\n\nimport (\n    \"fmt\"\n)\n\nfunc main() {\n    var X, Y int64\n    fmt.Scanf(\"%d %d\", &X, &Y)\n    const MOD int64 = 998244353\n    var ans = ((X % MOD) * (Y % MOD)) % MOD\n\n    for i := X / Y; i <= X; i++ {\n        if i == 0 {\n            continue\n        }\n        if i*i > X {\n            for j := int64(1); j <= Y; j++ {\n                if X/j >= i {\n                    ans -= ((X / j) % MOD) * (j % MOD) % MOD\n                    ans = (ans + MOD) % MOD\n                } else {\n                    break\n                }\n            }\n            break\n        } else {\n            M := min(X/i, Y)\n            m := X / (i + 1)\n            p := M + m + 1\n            q := M - m\n            if p%2 == 0 { p /= 2 }\n            if q%2 == 0 { q /= 2 }\n            r := (p % MOD) * (q % MOD) % MOD\n            r *= i\n            r %= MOD\n            ans -= r\n            ans = (ans + MOD) % MOD\n        }\n    }\n    if ans < 0 { ans += MOD }\n    fmt.Println(ans)\n}\n\nfunc min(a, b int64) int64 {\n    if a < b {\n        return a\n    }\n    return b\n}",
    "submitted_at": "2024-02-09T16:26:02Z"
}

// C++
{
  "solution_id": 8,
  "user_id": 10,
  "problem_id": 2,
  "language_id": 2,
  "code": "#include<bits/stdc++.h>\nusing namespace std;\n\nconst long long MOD = 998244353;\n\nsigned main(){\n    long long X,Y;\n    cin>>X>>Y;\n    long long ans = (X%MOD)*(Y%MOD)%MOD;\n\n    for (long long i = X/Y;i <=X;i++){\n        if (i == 0){continue;}\n        if (i > X / i){\n            for (long long j = 1;j <= Y;j++){\n                if (X/j >= i){\n                    ans -= ((X/j)%MOD) * (j%MOD)%MOD;\n                    ans %= MOD;\n                }\n                else{\n                    break;\n                }\n            }\n            break;\n        }\n        else{\n            long long M = min(X/i,Y);\n            long long m = X/(i + 1);\n            long long p = M + m + 1;\n            long long q = M - m;\n            if (p % 2 == 0){p /= 2;}\n            if (q % 2 == 0){q /= 2;}\n            long long r = (p % MOD) * (q % MOD)%MOD;\n            r *= i;\n            r %= MOD;\n            ans -= r;\n            ans %= MOD;\n        }\n    }\n    ans %= MOD;\n    if (ans < 0){ans += MOD;}\n    cout << ans << endl;\n    return 0;\n}",
  "submitted_at": "2024-02-09T16:26:02Z"
}

// python
{"solution_id":9,"user_id":10,"problem_id":2,"language_id":1,"code":"X, Y = map(int, input().split())\nMOD = 998244353\nMAX = int(1e12)\n\n\nassert 1 \u003c= X \u003c= MAX\nassert 1 \u003c= Y \u003c= MAX\n\nans = X * Y\n\nfor i in range(X // Y, X + 1):\n    if i * i \u003e X:\n        for j in range(1, Y + 1):\n            if X // j \u003e= i:\n                ans -= (X // j) * j\n            else:\n                break\n        break\n    if i == 0:\n        continue\n\n    M, m = min(X // i, Y), X // (i + 1)\n    ans -= (M * (M + 1) // 2 - m * (m + 1) // 2) * i\n\nans %= MOD\nprint(ans)","submitted_at":"2024-02-09T16:26:02Z"}
```