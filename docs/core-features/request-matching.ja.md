# リクエストマッチング

GoStubbyは、モックレスポンスを返すべき条件を正確に定義できる強力なリクエストマッチング機能を提供します。このドキュメントでは、利用可能な様々なリクエストマッチングオプションについて説明します。

## URLパステンプレート

GoStubbyは、URLパスのマッチングにテンプレートベースのシステムを使用します。テンプレートには固定セグメントと変数パラメータを含めることができます。

### 基本的なテンプレート構造

```
/fixed/path/{variable}/segments/{another_variable}
```

### 例

```json
{
  "urlPathTemplate": "/users/{id}/posts/{postId}",
  "method": "GET"
}
```

### 変数の命名規則
- 英数字とアンダースコアを使用可能
- 大文字と小文字は区別される
- テンプレート内で一意である必要がある
- 数字で始めることはできない

## マッチングタイプ

GoStubbyは、パスパラメータ、クエリパラメータ、リクエストボディに適用できる複数のマッチングパターンをサポートしています。

### 1. 完全一致 (`equalTo`)

指定された文字列と完全に一致する必要があります。

```json
{
  "pathParameters": {
    "id": {
      "equalTo": "12345"
    }
  }
}
```

### 2. 正規表現 (`matches`)

値を正規表現パターンと照合します。

```json
{
  "pathParameters": {
    "id": {
      "matches": "^[0-9]+$"
    }
  }
}
```

### 3. 否定的な正規表現 (`doesNotMatch`)

値が正規表現パターンに一致しないことを確認します。

```json
{
  "pathParameters": {
    "username": {
      "doesNotMatch": "[0-9]+"
    }
  }
}
```

### 4. 含む (`contains`)

値が指定された部分文字列を含むかどうかをチェックします。

```json
{
  "queryParameters": {
    "tags": {
      "contains": "important"
    }
  }
}
```

### 5. 含まない (`doesNotContain`)

値が指定された部分文字列を含まないことを確認します。

```json
{
  "queryParameters": {
    "status": {
      "doesNotContain": "deleted"
    }
  }
}
```

## パラメータタイプ

### パスパラメータ

URLパス変数のバリデーションルールを定義します。

```json
{
  "urlPathTemplate": "/users/{id}",
  "pathParameters": {
    "id": {
      "matches": "^[0-9]+$"
    }
  }
}
```

### クエリパラメータ

クエリ文字列パラメータを検証します。

```json
{
  "queryParameters": {
    "page": {
      "matches": "^[0-9]+$"
    },
    "limit": {
      "matches": "^[0-9]+$"
    },
    "sort": {
      "equalTo": "desc"
    }
  }
}
```

### リクエストボディ

リクエストボディの内容にマッチングパターンを適用します。

```json
{
  "body": {
    "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*"
  }
}
```

## 複数の条件

より正確な制御のために複数のマッチング条件を組み合わせることができます：

```json
{
  "request": {
    "urlPathTemplate": "/api/users/{id}",
    "method": "PUT",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    },
    "queryParameters": {
      "version": {
        "equalTo": "2.0"
      }
    },
    "body": {
      "matches": ".*\"status\":\\s*\"active\".*"
    }
  }
}
```

## ベストプラクティス

1. **パターンの具体性**
   - 意図しないマッチを避けるため、具体的なパターンを使用
   - パターンのエッジケースを考慮
   - 様々な入力でパターンをテスト

2. **正規表現**
   - パターンはシンプルで読みやすく保つ
   - 正規表現を徹底的にテスト
   - 複雑なパターンのパフォーマンスへの影響を考慮

3. **エラー処理**
   - マッチしないリクエストに適切なエラーレスポンスを提供
   - より良いエラーメッセージのための明確なバリデーションパターンを使用
   - バリデーション失敗時のカスタムエラーレスポンスを検討

4. **メンテナンス**
   - 複雑なパターンを文書化
   - 一貫した命名規則を使用
   - 関連するエンドポイントをグループ化

## 例

### 1. 基本的なAPIエンドポイント

```json
{
  "request": {
    "urlPathTemplate": "/api/v1/products/{id}",
    "method": "GET",
    "pathParameters": {
      "id": {
        "matches": "^[0-9]+$"
      }
    }
  }
}
```

### 2. クエリパラメータを使用した検索エンドポイント

```json
{
  "request": {
    "urlPathTemplate": "/api/v1/search",
    "method": "GET",
    "queryParameters": {
      "q": {
        "matches": ".{3,}"
      },
      "category": {
        "matches": "^(electronics|books|clothing)$"
      }
    }
  }
}
```

### 3. ボディバリデーションを使用したユーザー作成

```json
{
  "request": {
    "urlPathTemplate": "/api/v1/users",
    "method": "POST",
    "body": {
      "matches": ".*\"email\":\\s*\"[^@]+@[^@]+\\.[^@]+\".*",
      "contains": "\"age\":",
      "doesNotContain": "\"password\":"
    }
  }
}
```

## トラブルシューティング

1. **パターンが一致しない**
   - パターンの構文を確認
   - まずシンプルなパターンでテスト
   - 複雑なパターンの場合は正規表現テストツールを使用

2. **複数のマッチ**
   - パターンの具体性を見直し
   - 競合するパターンがないか確認
   - パターンの順序を考慮

3. **パフォーマンスの問題**
   - 複雑なパターンを簡素化
   - 複雑な正規表現の過度な使用を避ける
   - 頻繁に使用されるパターンのキャッシュを検討
