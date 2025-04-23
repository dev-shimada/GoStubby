# SSL/TLS 設定ガイド

このガイドでは、GoStubbyのSSL/TLS実装について、セットアップ、設定、およびモックサーバーのセキュア化のためのベストプラクティスを説明します。

## 概要

GoStubbyは、SSL/TLSを通じてHTTPSをサポートし、以下の機能を提供します：
- セキュアなHTTPSエンドポイントの実行
- HTTPとHTTPSの同時サポート
- カスタム証明書の設定
- 最新のセキュリティ標準の強制

## クイックスタート

### 基本的なHTTPS設定

1. 自己署名証明書の生成（開発用）：
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key -out ./certs/server.crt
```

2. HTTPSを有効にしてサーバーを起動：
```bash
go run main.go --cert ./certs/server.crt --key ./certs/server.key
```

## コマンドラインオプション

### SSL/TLS設定フラグ

| フラグ | 長形式 | 説明 | デフォルト |
|------|-----------|-------------|---------|
| `-s` | `--https-port` | HTTPSポート番号 | 8443 |
| `-t` | `--cert` | SSL/TLS証明書ファイルへのパス | - |
| `-k` | `--key` | SSL/TLS秘密鍵ファイルへのパス | - |

### 使用例

```bash
# カスタムHTTPSポートで実行
go run main.go --cert ./certs/server.crt --key ./certs/server.key --https-port 443

# HTTPとHTTPSを異なるポートで同時に実行
go run main.go --port 8080 --https-port 8443 \
  --cert ./certs/server.crt --key ./certs/server.key
```

## 証明書管理

### 1. 自己署名証明書

開発用の自己署名証明書の生成：

```bash
# 秘密鍵と証明書の生成
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key \
  -out ./certs/server.crt \
  -subj "/CN=localhost"

# カスタム詳細付きの生成
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ./certs/server.key \
  -out ./certs/server.crt \
  -subj "/C=JP/ST=Tokyo/L=渋谷区/O=組織名/CN=localhost"
```

### 2. Let's Encryptの使用

本番環境では、Let's Encryptから無料の証明書を取得：

1. certbotのインストール：
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install certbot

# macOS
brew install certbot
```

2. 証明書の生成：
```bash
sudo certbot certonly --standalone -d あなたのドメイン.com
```

3. 生成された証明書の使用：
```bash
go run main.go \
  --cert /etc/letsencrypt/live/あなたのドメイン.com/fullchain.pem \
  --key /etc/letsencrypt/live/あなたのドメイン.com/privkey.pem
```

### 3. 商用証明書

商用SSL証明書を使用する場合：

1. 証明書ファイルの結合：
```bash
cat domain.crt intermediate.crt root.crt > fullchain.pem
```

2. 結合した証明書でサーバーを起動：
```bash
go run main.go \
  --cert ./certs/fullchain.pem \
  --key ./certs/private.key
```

## セキュリティ設定

### TLSバージョン制御

GoStubbyは現代的なTLS標準を強制します：
- 最小TLSバージョン: 1.2
- 推奨される暗号スイート
- 完全前方秘匿性（PFS）

### 暗号スイートの設定

最大のセキュリティのために設定されるデフォルトの暗号スイート：

```go
[]uint16{
    tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
    tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
    tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
    tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}
```

## ベストプラクティス

### 1. 証明書管理
- 秘密鍵を安全に保管
- 定期的な証明書のローテーション
- 証明書の有効期限を監視
- 適切な鍵サイズの使用（RSAの場合最小2048ビット）

### 2. セキュリティヘッダー

レスポンスにセキュリティヘッダーを設定：

```json
{
  "response": {
    "headers": {
      "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
      "X-Content-Type-Options": "nosniff",
      "X-Frame-Options": "DENY",
      "X-XSS-Protection": "1; mode=block"
    }
  }
}
```

### 3. 本番環境設定
- 信頼できるCAからの有効なSSL証明書を使用
- HTTP/2サポートの有効化
- 適切なCORSヘッダーの設定
- レート制限の実装
- セキュリティログの監視

## トラブルシューティング

### よくある問題

1. **証明書の問題**
```
Error: tls: failed to load certificate
```
- ファイルパスの確認
- ファイルのパーミッションの確認
- 証明書フォーマットの確認
- 証明書と秘密鍵の対応を確認

2. **ポートアクセスの問題**
```
Error: listen tcp :443: bind: permission denied
```
- 非root利用者の場合は1024より大きいポートを使用
- 適切なシステムパーミッションを設定
- 必要に応じてポートフォワーディングを使用

3. **証明書信頼の問題**
```
Error: x509: certificate signed by unknown authority
```
- ルート証明書を信頼ストアに追加
- 適切な証明書チェーンを使用
- 中間証明書を確認

### 検証ツール

1. SSL Labs テスト：
```bash
# サーバーのテスト
curl https://www.ssllabs.com/ssltest/analyze.html?d=あなたのドメイン.com
```

2. OpenSSL検証：
```bash
# 証明書の検証
openssl x509 -in server.crt -text -noout

# 接続のテスト
openssl s_client -connect localhost:8443 -tls1_2
```

## セキュリティ考慮事項

1. **証明書の保管**
- 秘密鍵の安全な保管
- 適切なアクセス制御の実装
- 本番環境ではHSMの使用を検討

2. **更新とメンテナンス**
- TLSライブラリを最新に保つ
- セキュリティアドバイザリの監視
- 証明書の更新計画
- 定期的なセキュリティ監査

3. **ログとモニタリング**
- TLSハンドシェイクの失敗をログ記録
- 証明書の有効期限を監視
- セキュリティヘッダーの遵守を追跡
- セキュリティイベントのアラート設定
