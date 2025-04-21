# GoStubby
[![Go Report Card](https://goreportcard.com/badge/github.com/dev-shimada/GoStubby)](https://goreportcard.com/report/github.com/dev-shimada/GoStubby)
[![Coverage Status](https://coveralls.io/repos/github/dev-shimada/GoStubby/badge.svg?branch=main)](https://coveralls.io/github/dev-shimada/GoStubby?branch=main)
[![CI](https://github.com/dev-shimada/GoStubby/actions/workflows/CI.yaml/badge.svg)](https://github.com/dev-shimada/GoStubby/actions/workflows/CI.yaml)
[![build](https://github.com/dev-shimada/GoStubby/actions/workflows/build-docker-image.yaml/badge.svg)](https://github.com/dev-shimada/GoStubby/actions/workflows/build-docker-image.yaml)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/dev-shimada/GoStubby/blob/master/LICENSE)

柔軟で強力なGoによるモックサーバーの実装です。高度なリクエストマッチング機能とテンプレート化されたレスポンスを使用してモックエンドポイントを定義することができます。

## 特徴

- **柔軟なリクエストマッチング**:
  - テンプレートを使用したURLパスマッチング（例：`/users/{id}`）
  - 正規表現パターンマッチング
  - クエリパラメータのバリデーション
  - リクエストボディのバリデーション
  - 複数のマッチングパターン：`equalTo`、`matches`、`doesNotMatch`、`contains`、`doesNotContain`

- **強力なレスポンス処理**:
  - リクエストパラメータにアクセス可能なテンプレートベースのレスポンスボディ
  - ファイルベースのレスポンスボディ
  - カスタムHTTPステータスコード
  - カスタムレスポンスヘッダー

## インストール

```bash
go get github.com/dev-shimada/gostubby
```

## 使用方法

1. `./configs/config.json`に設定ファイルを作成：

```json
[
  {
    "request": {
      "urlPathTemplate": "/users/{id}",
      "method": "GET",
      "pathParameters": {
        "id": {
          "matches": "^[0-9]+$"
        }
      }
    },
    "response": {
      "status": 200,
      "body": "{\"id\": \"{{.Path.id}}\", \"name\": \"User {{.Path.id}}\"}"
    }
  }
]
```

2. サーバーを起動：

```bash
# デフォルトポート（8080）で起動
go run main.go

# 短いオプション（-p）でポートを指定
go run main.go -p 3000

# 長いオプション（--port）でポートを指定
go run main.go --port 3000
```

デフォルトでポート8080で起動します。`-p`または`--port`オプションで別のポート番号を指定することができます。

## 設定フォーマット

### リクエストマッチング

```json
{
  "request": {
    "urlPathTemplate": "/example/{param}",  // パスパラメータを含むURLテンプレート
    "method": "GET",                        // HTTPメソッド
    "pathParameters": {                     // パスパラメータのバリデーションルール
      "param": {
        "equalTo": "value",                 // 完全一致
        "matches": "^[0-9]+$",             // 正規表現パターンマッチ
        "doesNotMatch": "[a-z]+",          // 否定的な正規表現パターンマッチ
        "contains": "substring",            // 文字列を含む
        "doesNotContain": "substring"       // 文字列を含まない
      }
    },
    "queryParameters": {                    // クエリパラメータのバリデーション
      "param": {
        // パスパラメータと同じマッチングルール
      }
    },
    "body": {                              // リクエストボディのバリデーション
      // パラメータと同じマッチングルール
    }
  }
}
```

### レスポンス設定

```json
{
  "response": {
    "status": 200,                         // HTTPステータスコード
    "body": "Response content",            // 直接のレスポンス内容
    "bodyFileName": "response.json",       // または、ファイルベースのレスポンス
    "headers": {                           // カスタムレスポンスヘッダー
      "Content-Type": "application/json"
    }
  }
}
```

### テンプレート変数

レスポンスボディ内で以下のテンプレート変数を使用できます：
- パスパラメータ：`{{.Path.paramName}}`
- クエリパラメータ：`{{.Query.paramName}}`

## 設定例

1. パスパラメータを持つ基本的なエンドポイント：
```json
{
  "request": {
    "urlPathTemplate": "/users/{id}",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  },
  "response": {
    "status": 200,
    "body": "{\"id\": \"{{.Path.id}}\", \"name\": \"User {{.Path.id}}\"}"
  }
}
```

2. ファイルベースのレスポンスを持つエンドポイント：
```json
{
  "request": {
    "urlPathTemplate": "/data/{type}",
    "method": "GET"
  },
  "response": {
    "status": 200,
    "bodyFileName": "responses/data.json"
  }
}
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています - 詳細は[LICENSE](LICENSE)ファイルを参照してください。

*他の言語で読む: [English](README.md)*
