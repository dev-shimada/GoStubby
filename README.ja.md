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
# デフォルトポート（8080）とデフォルト設定ディレクトリで起動
go run main.go

# 短いオプション（-p）でポートを指定
go run main.go -p 3000

# 長いオプション（--port）でポートを指定
go run main.go --port 3000

# 短いオプション（-c）で設定ファイルのパスを指定
go run main.go -c ./path/to/config.json

# 長いオプション（--config）で設定ディレクトリを指定
go run main.go --config ./path/to/configs
```

サーバーは以下のコマンドラインオプションをサポートしています：

HTTP設定：
- ポート番号: `-p` または `--port`（デフォルト: 8080）

HTTPS設定：
- HTTPSポート番号: `-s` または `--https-port`（デフォルト: 8443）
- SSL/TLS証明書: `-t` または `--cert`（SSL/TLS証明書ファイルへのパス）
- SSL/TLS秘密鍵: `-k` または `--key`（SSL/TLS秘密鍵ファイルへのパス）

一般設定：
- 設定ファイル: `-c` または `--config`（デフォルト: "./configs"）

設定ファイルは、単一のJSONファイルまたは複数のJSONファイルを含むディレクトリのいずれかを指定できます。ディレクトリを指定した場合、そのディレクトリ内のすべてのJSONファイルが読み込まれます。

### SSL/TLSサポート

SSL/TLS証明書を提供することで、サーバーをHTTPSモードで実行できます。HTTPとHTTPSを同時に有効にして実行することも可能です。

HTTPSを有効にするには：
1. SSL/TLS証明書と秘密鍵ファイルを用意
2. 証明書と鍵ファイルのパスを指定してサーバーを起動：

```bash
# HTTPとHTTPSの両方で実行
go run main.go --cert ./certs/server.crt --key ./certs/server.key

# HTTPとHTTPSのポートをカスタマイズ
go run main.go --port 8080 --https-port 8443 --cert ./certs/server.crt --key ./certs/server.key
```

開発やテスト用に自己署名証明書を生成する場合：
```bash
# 秘密鍵と自己署名証明書の生成
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key -out ./certs/server.crt
```

注意：セキュリティのため、TLS 1.2以上のバージョンを強制しています。

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
