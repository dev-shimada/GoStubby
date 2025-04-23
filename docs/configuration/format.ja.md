# 設定フォーマットガイド

このドキュメントでは、GoStubbyの設定フォーマットについて、利用可能なすべてのオプションと例を含めて詳しく説明します。

## 設定構造

GoStubbyは設定ファイルにJSONフォーマットを使用します。各設定ファイルにはスタブマッピングの配列が含まれます：

```json
[
  {
    "request": {
      // リクエストマッチングの設定
    },
    "response": {
      // レスポンスの設定
    }
  }
]
```

## 完全な設定スキーマ

```json
{
  "request": {
    "urlPathTemplate": string,          // パスパラメータを含むURLテンプレート
    "method": string,                   // HTTPメソッド（GET, POST等）
    "pathParameters": {                 // パスパラメータのバリデーションルール
      "paramName": {
        "equalTo": string,             // 完全一致
        "matches": string,             // 正規表現パターン
        "doesNotMatch": string,        // 否定的な正規表現パターン
        "contains": string,            // 文字列を含む
        "doesNotContain": string       // 文字列を含まない
      }
    },
    "queryParameters": {               // クエリパラメータのバリデーション
      "paramName": {
        // パスパラメータと同じルール
      }
    },
    "body": {                         // リクエストボディのバリデーション
      // パラメータと同じルール
    }
  },
  "response": {
    "status": number,                 // HTTPステータスコード
    "body": string,                   // 直接のレスポンス内容
    "bodyFileName": string,           // ファイルベースのレスポンス
    "headers": {                      // レスポンスヘッダー
      "headerName": string
    }
  }
}
```

## リクエスト設定

### URLパステンプレート

テンプレートは固定セグメントと変数パラメータをサポートします：

```json
{
  "urlPathTemplate": "/api/v1/users/{id}/posts/{postId}"
}
```

### HTTPメソッド

サポートされるHTTPメソッド：
- GET
- POST
- PUT
- DELETE
- PATCH
- HEAD
- OPTIONS

```json
{
  "method": "POST"
}
```

### パラメータバリデーション

#### パスパラメータ

```json
{
  "pathParameters": {
    "id": {
      "matches": "^[0-9]+$"
    },
    "category": {
      "equalTo": "electronics"
    }
  }
}
```

#### クエリパラメータ

```json
{
  "queryParameters": {
    "page": {
      "matches": "^[0-9]+$"
    },
    "sort": {
      "equalTo": "desc"
    }
  }
}
```

#### リクエストボディ

```json
{
  "body": {
    "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*",
    "contains": "\"active\":true"
  }
}
```

## レスポンス設定

### ステータスコード

```json
{
  "status": 200  // 有効なHTTPステータスコード
}
```

### 直接レスポンスボディ

```json
{
  "body": "{\"message\": \"成功\"}"
}
```

### ファイルベースのレスポンス

```json
{
  "bodyFileName": "responses/user-profile.json"
}
```

### カスタムヘッダー

```json
{
  "headers": {
    "Content-Type": "application/json",
    "Cache-Control": "no-cache",
    "X-Custom-Header": "カスタム値"
  }
}
```

## テンプレート変数

レスポンステンプレートで使用可能な変数：

```json
{
  "body": {
    "path": "{{.Path.paramName}}",      // パスパラメータ
    "query": "{{.Query.paramName}}",     // クエリパラメータ
    "method": "{{.Request.Method}}",     // HTTPメソッド
    "header": "{{.Request.Header.name}}" // リクエストヘッダー
  }
}
```

## 設定管理

### ファイル構成

推奨されるディレクトリ構造：
```
configs/
├── api/
│   ├── users.json
│   └── products.json
├── mock/
│   └── test-data.json
└── config.json
```

### 複数の設定ファイル

複数のファイルを使用する場合：
1. 各ファイルは有効なJSON配列を含む必要があります
2. ファイルはアルファベット順に読み込まれます
3. 後の定義が先の定義を上書きします

## 例

### 1. 基本的なRESTエンドポイント

```json
{
  "request": {
    "urlPathTemplate": "/api/users/{id}",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  },
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "body": "{\"id\": {{.Path.id}}, \"name\": \"ユーザー {{.Path.id}}\"}"
  }
}
```

### 2. 複雑なバリデーション

```json
{
  "request": {
    "urlPathTemplate": "/api/products",
    "method": "POST",
    "queryParameters": {
      "version": {
        "equalTo": "2.0"
      }
    },
    "body": {
      "matches": ".*\"price\":\\s*[0-9]+(\\.?[0-9]*)?.*",
      "contains": "\"category\":"
    }
  },
  "response": {
    "status": 201,
    "headers": {
      "Content-Type": "application/json",
      "Location": "/api/products/{{.Response.id}}"
    },
    "body": "{\"message\": \"商品が作成されました\", \"id\": \"12345\"}"
  }
}
```

### 3. ファイルベースの設定

```json
{
  "request": {
    "urlPathTemplate": "/api/data/{type}",
    "method": "GET",
    "pathParameters": {
      "type": {
        "matches": "^(users|products|orders)$"
      }
    }
  },
  "response": {
    "status": 200,
    "headers": {
      "Content-Type": "application/json"
    },
    "bodyFileName": "responses/{{.Path.type}}.json"
  }
}
```

## ベストプラクティス

1. **ファイル構成**
   - 意味のあるファイル名を使用
   - 関連するスタブをグループ化
   - 一貫した構造を維持

2. **バリデーションルール**
   - 正規表現パターンをシンプルに
   - 具体的なマッチングルールを使用
   - 複雑なパターンを文書化

3. **レスポンス管理**
   - 動的コンテンツにテンプレートを使用
   - レスポンスファイルを論理的に整理
   - 一貫したレスポンスフォーマットを維持

4. **バージョン管理**
   - 設定のバージョン管理
   - 変更を文書化
   - 意味のあるコミットメッセージを使用

## トラブルシューティング

1. **不正なJSON**
   - JSONバリデータを使用
   - 構文エラーを確認
   - ファイルエンコーディングを確認

2. **パターンマッチングの問題**
   - 正規表現パターンをテスト
   - URLテンプレートを確認
   - パラメータ名を確認

3. **ファイル読み込み**
   - ファイルパスを確認
   - パーミッションを確認
   - JSON構造を検証
