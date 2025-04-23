# レスポンス処理

GoStubbyは、静的および動的なモックレスポンスを作成できる柔軟なレスポンス処理機能を提供します。このドキュメントでは、レスポンスの設定とカスタマイズのすべての側面について説明します。

## レスポンス設定

### 基本構造

```json
{
  "response": {
    "status": 200,                         // HTTPステータスコード
    "body": "レスポンス内容",               // 直接のレスポンス内容
    "bodyFileName": "response.json",       // またはファイルベースのレスポンス
    "headers": {                           // カスタムレスポンスヘッダー
      "Content-Type": "application/json"
    }
  }
}
```

## レスポンスタイプ

### 1. 直接レスポンスボディ

設定内で直接レスポンス内容を指定するには `body` フィールドを使用します：

```json
{
  "response": {
    "status": 200,
    "body": "{\"message\": \"こんにちは、世界！\"}",
    "headers": {
      "Content-Type": "application/json"
    }
  }
}
```

### 2. ファイルベースのレスポンス

ファイルからレスポンス内容を読み込むには `bodyFileName` フィールドを使用します：

```json
{
  "response": {
    "status": 200,
    "bodyFileName": "responses/user-profile.json",
    "headers": {
      "Content-Type": "application/json"
    }
  }
}
```

## テンプレートベースのレスポンス

GoStubbyは、テンプレートを使用した動的なレスポンス生成をサポートしています。テンプレートはリクエストパラメータにアクセスし、カスタマイズされたレスポンスを生成できます。

### 利用可能なテンプレート変数

1. **パスパラメータ**
   ```json
   {
     "response": {
       "body": "{\"userId\": \"{{.Path.id}}\", \"message\": \"ユーザー{{.Path.id}}の詳細\"}"
     }
   }
   ```

2. **クエリパラメータ**
   ```json
   {
     "response": {
       "body": "{\"search\": \"{{.Query.q}}\", \"page\": {{.Query.page}}}"
     }
   }
   ```

### テンプレート構文

- パスパラメータ: `{{.Path.paramName}}`
- クエリパラメータ: `{{.Query.paramName}}`
- HTTPメソッド: `{{.Request.Method}}`
- リクエストヘッダー: `{{.Request.Header.headerName}}`

## ステータスコード

異なるシナリオに適切なHTTPステータスコードを設定します：

```json
{
  "response": {
    "status": 201,  // 作成成功
    "body": "{\"message\": \"リソースが正常に作成されました\"}"
  }
}
```

一般的なステータスコード：
- 200: OK（成功）
- 201: Created（作成成功）
- 400: Bad Request（不正なリクエスト）
- 401: Unauthorized（未認証）
- 403: Forbidden（アクセス禁止）
- 404: Not Found（見つからない）
- 500: Internal Server Error（サーバーエラー）

## レスポンスヘッダー

### カスタムヘッダーの設定

```json
{
  "response": {
    "headers": {
      "Content-Type": "application/json",
      "Cache-Control": "no-cache",
      "X-Custom-Header": "カスタム値"
    }
  }
}
```

### 一般的なヘッダー
- Content-Type
- Cache-Control
- Access-Control-Allow-Origin
- Authorization
- X-Rate-Limit

## 高度なレスポンス機能

### 1. 条件付きレスポンス

リクエストパラメータに基づいて異なるテンプレートを使用：

```json
{
  "response": {
    "body": "{% if eq .Query.type \"premium\" %}
      {\"message\": \"プレミアムコンテンツ\"}
    {% else %}
      {\"message\": \"基本コンテンツ\"}
    {% endif %}"
  }
}
```

### 2. 動的ファイル読み込み

パラメータに基づいて異なるレスポンスファイルを読み込み：

```json
{
  "response": {
    "bodyFileName": "responses/{{.Path.type}}.json"
  }
}
```

### 3. カスタムエラーレスポンス

```json
{
  "response": {
    "status": 400,
    "body": "{\"error\": \"不正なリクエスト\", \"details\": \"必須フィールドがありません: {{.Path.field}}\"}"
  }
}
```

## ベストプラクティス

1. **レスポンスの整理**
   - 関連するレスポンスをディレクトリでグループ化
   - 意味のあるファイル名を使用
   - 一貫したファイル構造を維持

2. **テンプレートの使用**
   - テンプレートはシンプルで読みやすく
   - テンプレート構文を検証
   - 欠落したパラメータを適切に処理

3. **エラー処理**
   - 適切なステータスコードを使用
   - 意味のあるエラーメッセージを提供
   - 関連するエラー詳細を含める

4. **パフォーマンス**
   - 可能な場合はファイルベースのレスポンスをキャッシュ
   - テンプレートの複雑さを最小限に
   - 適切なコンテンツ圧縮を使用

## 例

### 1. 基本的なJSONレスポンス

```json
{
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"id\": 1, \"name\": \"例\"}"
  }
}
```

### 2. テンプレートを使用した動的レスポンス

```json
{
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{
      \"userId\": \"{{.Path.id}}\",
      \"query\": \"{{.Query.q}}\",
      \"timestamp\": \"{{.Request.Time}}\"
    }"
  }
}
```

### 3. カスタムヘッダー付きファイルベースレスポンス

```json
{
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json",
      "Cache-Control": "max-age=3600",
      "ETag": "\"123456\""
    },
    "bodyFileName": "responses/large-response.json"
  }
}
```

## トラブルシューティング

1. **テンプレートエラー**
   - テンプレート構文を確認
   - パラメータ名を確認
   - パラメータが利用可能か確認

2. **ファイル読み込みの問題**
   - ファイルパスを確認
   - ファイルのパーミッションを確認
   - ファイル内容を検証

3. **コンテンツタイプの不一致**
   - Content-Typeヘッダーがボディと一致しているか確認
   - JSON構文を確認
   - 文字エンコーディングを確認

4. **パフォーマンスの問題**
   - ファイルサイズを監視
   - テンプレート処理を最適化
   - レスポンスのキャッシュを検討
