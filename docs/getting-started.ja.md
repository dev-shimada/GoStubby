# GoStubby 入門ガイド

このガイドでは、GoStubbyを素早くセットアップして使い始めるための手順を説明します。

## インストール

GoStubbyをインストールして使用するには、以下の2つの方法があります：

### 1. Goパッケージのインストール

Go のパッケージ管理を使用してインストール：

```bash
go install github.com/dev-shimada/gostubby@latest
```

### 2. Dockerイメージ

公式Dockerイメージを使用：

```bash
# イメージの取得
docker pull ghcr.io/dev-shimada/gostubby:latest

# コンテナの実行
docker run -p 8080:8080 -v $(pwd)/configs:/app/configs ghcr.io/dev-shimada/gostubby:latest
```

Dockerを使用する利点：
- Goのインストールが不要
- プラットフォーム間で一貫した環境
- コンテナ化環境での容易なデプロイ
- 最新イメージの取得による自動アップデート

## クイックスタート チュートリアル

### 1. 最初の設定ファイルの作成

`configs/config.json` という新しいファイルを作成し、基本的なエンドポイント設定を記述します：

```json
[
  {
    "request": {
      "urlPathTemplate": "/hello/{name}",
      "method": "GET",
      "pathParameters": {
        "name": {
          "matches": "^[a-zA-Z]+$"
        }
      }
    },
    "response": {
      "status": 200,
      "headers": {
        "Content-Type": "application/json"
      },
      "body": "{\"message\": \"こんにちは、{{.Path.name}}さん！\"}"
    }
  }
]
```

### 2. サーバーの起動

デフォルト設定でサーバーを起動します：

```bash
go run main.go
```

サーバーは以下のポートで起動します：
- HTTP: `http://localhost:8080`
- HTTPS（設定されている場合）: `https://localhost:8443`

### 3. モックエンドポイントのテスト

curl または任意の HTTP クライアントを使用してエンドポイントをテストします：

```bash
curl -s -X GET "http://localhost:8080/example/file?param1=123%3F" | iconv -f sjis -t utf8
```

期待されるレスポンス：
```json
{
    "message": "This is a stub response", 
    "description": "テスト",
    "path1": "file",
    "param1": 123?
}
```

## 基本的な設定例

### 1. 静的レスポンス

```json
{
  "request": {
    "urlPathTemplate": "/api/status",
    "method": "GET"
  },
  "response": {
    "status": 200,
    "body": "{\"status\": \"正常稼働中\"}"
  }
}
```

### 2. パスパラメータを使用した動的レスポンス

```json
{
  "request": {
    "urlPathTemplate": "/users/{id}/profile",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  },
  "response": {
    "status": 200,
    "body": "{\"userId\": \"{{.Path.id}}\", \"name\": \"ユーザー{{.Path.id}}\"}"
  }
}
```

### 3. リクエストボディのバリデーション

```json
{
  "request": {
    "urlPathTemplate": "/api/users",
    "method": "POST",
    "body": {
      "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*"
    }
  },
  "response": {
    "status": 201,
    "body": "{\"message\": \"ユーザーが正常に作成されました\"}"
  }
}
```

## コマンドラインオプション

GoStubby は以下のコマンドラインオプションをサポートしています：

### HTTP設定
- ポート番号: `-p` または `--port`（デフォルト: 8080）

### HTTPS設定
- HTTPSポート番号: `-s` または `--https-port`（デフォルト: 8443）
- SSL/TLS証明書: `-t` または `--cert`（SSL/TLS証明書ファイルへのパス）
- SSL/TLS秘密鍵: `-k` または `--key`（SSL/TLS秘密鍵ファイルへのパス）

### 一般設定
- 設定ファイル: `-c` または `--config`（デフォルト: "./configs"）

カスタム設定の例：
```bash
go run main.go --port 3000 --config ./my-configs
```

## HTTPSの設定

1. 自己署名証明書の生成（開発用）：
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key -out ./certs/server.crt
```

2. HTTPSでサーバーを起動：
```bash
go run main.go --cert ./certs/server.crt --key ./certs/server.key
```

## 次のステップ

- [リクエストマッチング](core-features/request-matching.ja.md)で高度なリクエストマッチングについて学ぶ
- [レスポンス処理](core-features/response-handling.ja.md)で複雑なレスポンスについて学ぶ
- [SSL/TLS](security/ssl-tls.ja.md)でセキュアなエンドポイントを設定
- [設定ガイド](configuration/format.ja.md)で詳細な設定オプションを確認

## よくある問題と解決方法

1. **ポートが既に使用されている**
   ```bash
   # ポートを変更する
   go run main.go --port 3000
   ```

2. **設定ファイルが見つからない**
   - 設定ファイルが指定されたパスに存在することを確認
   - 絶対パスまたは正しい相対パスを使用
   - ファイルのパーミッションを確認

3. **設定フォーマットが無効**
   - JSONの構文を確認
   - 必須フィールドがすべて存在することを確認
   - マッチングパターンが正しいことを確認

4. **証明書の問題**
   - 証明書と鍵ファイルのパスを確認
   - ファイルが読み取り可能であることを確認
   - 証明書の有効期限と有効性を確認
