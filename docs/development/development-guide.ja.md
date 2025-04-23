# 開発ガイド

このガイドは、GoStubbyに貢献したい開発者向けの情報を提供します。開発環境のセットアップ、テストの実行、および貢献のガイドラインについて説明します。

## 開発環境のセットアップ

### 前提条件

1. Go 1.16以降
2. Git
3. Make（オプション、ただし推奨）
4. OpenSSL（SSL/TLS証明書生成用）

### 始め方

1. リポジトリのクローン：
```bash
git clone https://github.com/dev-shimada/GoStubby.git
cd GoStubby
```

2. 依存関係のインストール：
```bash
go mod download
```

3. プロジェクトのビルド：
```bash
go build
```

## プロジェクト構造

```
GoStubby/
├── .github/            # GitHub Actionsワークフロー
├── configs/            # 設定例
├── docs/              # ドキュメント
├── testdata/          # テストフィクスチャ
├── body/              # レスポンスボディテンプレート
├── main.go            # アプリケーションのエントリーポイント
├── main_test.go       # メインパッケージのテスト
├── go.mod             # Goモジュールファイル
└── go.sum             # Goモジュールチェックサム
```

## テスト

### テストの実行

すべてのテストを実行：
```bash
go test ./...
```

カバレッジ付きでテストを実行：
```bash
go test -cover ./...
```

カバレッジレポートの生成：
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### テストのカテゴリ

1. **ユニットテスト**
   - ソースファイルと同じディレクトリに配置
   - 名前のパターン: `*_test.go`
   - 個別のコンポーネントに焦点

2. **統合テスト**
   - `test`ディレクトリに配置
   - コンポーネント間の相互作用をテスト
   - テストフィクスチャを使用

3. **パフォーマンステスト**
   - 重要な操作のベンチマーク
   - `*_test.go`ファイルに配置
   - `testing.B`ベンチマークを使用

### テストの作成

テスト構造の例：
```go
func TestFeature(t *testing.T) {
    // テストのセットアップ
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "有効な入力",
            input:    "test",
            expected: "result",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テストの実行
            result := Feature(tt.input)
            
            // アサーション
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## コードスタイル

### Goのガイドライン

1. 標準のGo形式に従う：
```bash
gofmt -s -w .
```

2. golintの使用：
```bash
golint ./...
```

3. go vetの実行：
```bash
go vet ./...
```

### コードの構成

1. **パッケージ構造**
   - ディレクトリごとに1つのパッケージ
   - 明確なパッケージの責任
   - 最小限のパブリックインターフェース

2. **ファイルの構成**
   - 関連する機能をまとめる
   - 明確なファイル命名
   - 論理的なグループ化

3. **コードドキュメント**
   - エクスポートされるすべてのシンボルを文書化
   - 例を含める
   - 明確で簡潔なコメント

## 貢献

### 開発ワークフロー

1. **イシューの作成**
   - 問題/機能を説明
   - 関連するラベルを追加
   - 関連するイシューをリンク

2. **ブランチの作成**
   - 説明的な名前を使用
   - イシュー番号を含める
   - 例：`feature/123-add-ssl-support`

3. **開発**
   - 最初にテストを作成
   - 変更を実装
   - ドキュメントを更新

4. **コードレビュー**
   - プルリクエストを提出
   - レビューコメントに対応
   - 必要に応じて更新

### プルリクエストのガイドライン

1. **準備**
   - メインブランチにリベース
   - すべてのテストを実行
   - ドキュメントを更新

2. **PRの説明**
   - 明確な説明
   - イシューの参照
   - 変更点のリスト

3. **コード品質**
   - すべてのテストをパス
   - スタイルガイドに従う
   - ドキュメントを含める

### コミットメッセージ

conventional commitsフォーマットに従う：
```
type(scope): description

[optional body]

[optional footer]
```

タイプ：
- feat: 新機能
- fix: バグ修正
- docs: ドキュメント
- style: フォーマット
- refactor: コードのリストラクチャリング
- test: テストの追加
- chore: メンテナンス

例：
```
feat(ssl): SSL/TLSサポートの追加

- 証明書の読み込み機能追加
- HTTPSサーバーの実装
- ドキュメントの更新

Closes #123
```

## デバッグ

### VSCode設定

1. `.vscode/launch.json`の作成：
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "GoStubbyのデバッグ",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "args": ["--config", "./configs/debug.json"]
        }
    ]
}
```

2. ブレークポイントを追加
3. VSCodeのデバッグ機能を使用

### ロギング

1. 開発時のロギング：
```go
log.Printf("デバッグ: %v", value)
```

2. 本番環境のロギング：
```go
log.Printf("エラー: %v", err)
```

## パフォーマンスプロファイリング

### CPUプロファイリング

```bash
go test -cpuprofile cpu.prof -bench .
go tool pprof cpu.prof
```

### メモリプロファイリング

```bash
go test -memprofile mem.prof -bench .
go tool pprof mem.prof
```

## リリースプロセス

1. **バージョン更新**
   - バージョン番号の更新
   - CHANGELOG.mdの更新
   - ドキュメントの更新

2. **テスト**
   - すべてのテストを実行
   - 統合テストの実施
   - ドキュメントの確認

3. **リリース**
   - リリースブランチの作成
   - バージョンのタグ付け
   - リポジトリへのプッシュ

4. **リリース後**
   - メインブランチの更新
   - ブランチのクリーンアップ
   - ドキュメントの更新

## サポート

- GitHub Issues: バグ報告と機能リクエスト
- Discussions: 一般的な質問とディスカッション
- Pull Requests: コード貢献

注意事項：
- 既存のイシューを検索
- 明確な説明を提供
- 最小限の例を含める
- 敬意を持って建設的に
