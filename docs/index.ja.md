# GoStubby ドキュメント

GoStubbyのドキュメントへようこそ。このガイドでは、柔軟で強力なGoによるモックサーバーの実装であるGoStubbyの使用方法と貢献方法について詳しく説明します。

## ドキュメントセクション

### 📚 はじめに
- [プロジェクト概要](overview.md) - GoStubbyの紹介と機能
- [入門ガイド](getting-started.md) - インストールと基本的な使用方法
- [クイックスタート チュートリアル](getting-started.md#quick-start-tutorial) - すぐに始める方法

### 🔧 主要機能
- [リクエストマッチング](core-features/request-matching.md) - URLテンプレートとリクエストバリデーション
- [レスポンス処理](core-features/response-handling.md) - モックレスポンスの設定とカスタマイズ
- [テンプレートシステム](core-features/response-handling.md#template-based-responses) - 動的レスポンステンプレートの使用

### ⚙️ 設定
- [設定フォーマット](configuration/format.md) - 詳細な設定オプション
- [コマンドラインオプション](configuration/format.md#configuration-management) - 利用可能なCLI引数
- [設定例](configuration/format.md#examples) - 一般的な設定パターン

### 🔒 セキュリティ
- [SSL/TLS設定](security/ssl-tls.md) - モックサーバーのセキュア化
- [証明書管理](security/ssl-tls.md#certificate-management) - SSL証明書の取り扱い
- [セキュリティベストプラクティス](security/ssl-tls.md#best-practices) - 推奨されるセキュリティ設定

### 👩‍💻 開発
- [開発ガイド](development/development-guide.md) - GoStubbyへの貢献
- [テスト](development/development-guide.md#testing) - テストの実行と作成
- [コードスタイル](development/development-guide.md#code-style) - コーディング規約とガイドライン

## クイックナビゲーション

```mermaid
graph LR
    A[スタート] --> B[概要]
    A --> C[はじめに]
    
    B --> D[主要機能]
    C --> D
    
    D --> E[リクエストマッチング]
    D --> F[レスポンス処理]
    
    G[設定] --> H[フォーマット]
    G --> I[例]
    
    J[セキュリティ] --> K[SSL/TLS]
    J --> L[ベストプラクティス]
    
    M[開発] --> N[ガイド]
    M --> O[テスト]
```

## 機能概要

| 機能 | 説明 | ドキュメント |
|---------|-------------|---------------|
| URLテンプレート | 動的URLパスマッチング | [リクエストマッチング](core-features/request-matching.md#url-path-templates) |
| レスポンステンプレート | テンプレートベースのレスポンス | [レスポンス処理](core-features/response-handling.md#template-based-responses) |
| SSL/TLSサポート | セキュアなHTTPSエンドポイント | [SSL/TLS設定](security/ssl-tls.md) |
| 複数設定 | 柔軟な設定管理 | [設定フォーマット](configuration/format.md#configuration-management) |

## 一般的なタスク

### モックAPIのセットアップ
1. [GoStubbyのインストール](getting-started.md#installation)
2. [設定ファイルの作成](getting-started.md#quick-start-tutorial)
3. [サーバーの起動](getting-started.md#start-the-server)

### サーバーのセキュア化
1. [証明書の生成](security/ssl-tls.md#certificate-management)
2. [HTTPSの設定](security/ssl-tls.md#quick-start)
3. [セキュリティベストプラクティスの適用](security/ssl-tls.md#best-practices)

### コントリビューション
1. [開発環境のセットアップ](development/development-guide.md#development-environment-setup)
2. [コーディングガイドラインの確認](development/development-guide.md#code-style)
3. [プルリクエストの提出](development/development-guide.md#contributing)

## ヘルプの取得

- **Issues**: バグの報告や機能リクエストは[GitHub Issues](https://github.com/dev-shimada/GoStubby/issues)から
- **Discussions**: 質問は[GitHub Discussions](https://github.com/dev-shimada/GoStubby/discussions)で
- **貢献**: [開発ガイド](development/development-guide.md)を読んで貢献を開始

## バージョン情報

- 現在のバージョン: 1.0.0
- Go バージョン: ≥ 1.16
- ライセンス: MIT

最新の更新については、[GitHubリポジトリ](https://github.com/dev-shimada/GoStubby)をご確認ください。

---

*他の言語でこのドキュメントを読む: [English](index.md)*
