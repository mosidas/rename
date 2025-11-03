# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) およびAI開発者がこのリポジトリで作業する際のガイダンスを提供します。

## プロジェクト概要

**macOS専用ファイル一括リネームツール**

- **技術スタック**: Wails v2 + Go 1.21 + Next.js 15 (React 19) + TypeScript + Tailwind CSS 4
- **アーキテクチャ**: Clean Architecture / Hexagonal Architecture
- **開発手法**: TDD (Test-Driven Development)
- **設計原則**: SOLID原則

## アーキテクチャ詳細

### ディレクトリ構造

```
rename/
├── internal/                    # バックエンド（Go）
│   ├── domain/                 # ドメイン層
│   │   ├── strategy.go        # RenameStrategy interface + 実装
│   │   ├── file.go            # File entity
│   │   └── history.go         # History entity
│   ├── usecase/               # ユースケース層
│   │   ├── rename_usecase.go
│   │   └── history_usecase.go
│   ├── repository/            # リポジトリ層
│   │   └── json_history_repository.go
│   └── service/               # サービス層
│       └── filesystem_service.go
├── frontend/                   # フロントエンド（Next.js）
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx     # ルートレイアウト + Header
│   │   │   ├── page.tsx       # ホームページ
│   │   │   └── globals.css    # CSS変数（テーマ）
│   │   └── components/
│   │       ├── Header.tsx     # テーマ切り替え
│   │       └── RenamePanel.tsx # メインUI
│   └── wailsjs/               # Wails自動生成バインディング
├── app.go                      # プレゼンテーション層（Wails binding）
└── main.go                     # エントリーポイント
```

### レイヤー説明

#### 1. ドメイン層 (`internal/domain/`)

**責務**: ビジネスロジックとエンティティ

- **`strategy.go`**: Strategy Pattern実装
  - `RenameStrategy` interface
  - `ExactMatchStrategy`: 部分一致置換
  - `RegexMatchStrategy`: 正規表現置換
  - `CaseInsensitiveStrategy`: Decoratorパターンで大小文字無視

- **`file.go`**: ファイルエンティティ
  - 元のパス、新しいパス、変更検出

- **`history.go`**: 履歴エンティティ
  - 最大100件、重複時は最前面に移動

**重要**: ドメイン層は他のレイヤーに依存しない

#### 2. ユースケース層 (`internal/usecase/`)

**責務**: アプリケーションロジック

- **`RenameUseCase`**:
  - `GeneratePreview()`: プレビュー生成（Strategyを適用）
  - `Execute()`: 実際のリネーム実行
  - エラー時はスキップして続行
  - 成功/失敗件数、新ファイルパスを返却

- **`HistoryUseCase`**:
  - `GetHistory()`: 履歴取得
  - `AddEntry()`: 履歴追加（重複チェック）

**依存**: インターフェース（`FileSystemService`, `HistoryRepository`）に依存

#### 3. リポジトリ層 (`internal/repository/`)

**責務**: データ永続化

- **`JSONHistoryRepository`**:
  - 保存先: `~/.config/rename/config.json`
  - JSON形式で履歴を保存/読み込み

#### 4. サービス層 (`internal/service/`)

**責務**: 外部システムとの連携

- **`FileSystemService`**:
  - `RenameFile()`: os.Renameのラッパー
  - `FileExists()`: ファイル存在チェック

#### 5. プレゼンテーション層 (`app.go`)

**責務**: Wails bindingの薄いアダプター

- `SelectFiles()`: ファイル選択ダイアログ
- `GeneratePreview()`: プレビュー生成（UseCaseに委譲）
- `ExecuteRename()`: リネーム実行（UseCaseに委譲）
- `GetHistory()`: 履歴取得（UseCaseに委譲）

**重要**: ビジネスロジックを含まず、UseCaseに処理を委譲

### フロントエンド設計

#### テーマシステム

- **CSS変数ベース**: `globals.css`で定義
- **テーマ**: ライト/ダーク/システム
- **実装**: `data-theme`属性で切り替え
- **変数例**:
  ```css
  --background, --foreground, --border, --muted,
  --accent, --destructive, --success
  ```

#### RenamePanel.tsx（メインUI）

**主な機能**:
1. **リアルタイムプレビュー**: 300msデバウンス
2. **履歴ドロップダウン**: フォーカス時に表示、最大10件
3. **状態管理**: リネーム後もファイルと入力を保持
4. **レイアウト**: 1:2比率（入力欄:プレビュー）

**重要な実装**:
- `generatePreviewDebounced`: クロージャでdebounce実装
- `useEffect`: ファイル選択時に即座にプレビュー生成
- メッセージ表示: ボタン横に統合、自動クリアしない

#### Header.tsx

- テーマ切り替えボタン
- `localStorage`にテーマ保存
- システム設定との連携

## SOLID原則の適用

### Single Responsibility Principle (SRP)
- 各クラス/関数は単一の責務のみ
- 例: `RenameUseCase`はリネームロジックのみ、永続化は`Repository`

### Open/Closed Principle (OCP)
- Strategy Patternで拡張に開いている
- 新しい置換戦略は`RenameStrategy`を実装するだけ

### Liskov Substitution Principle (LSP)
- 全てのStrategyは`RenameStrategy`と互換

### Interface Segregation Principle (ISP)
- `FileSystemService`は最小限のメソッドのみ
- クライアントが不要なメソッドに依存しない

### Dependency Inversion Principle (DIP)
- UseCaseはインターフェースに依存
- 具象クラスは`NewApp()`で注入

## 開発ワークフロー

### TDD（Test-Driven Development）

**必須プロセス**:
1. **Red**: テストを先に書く（失敗することを確認）
2. **Green**: 最小限の実装で通す
3. **Refactor**: コードをクリーンに

**テストツール**: testify (assert, mock)

**テストコマンド**:
```bash
# 全テスト実行
go test ./... -v

# カバレッジ付き
go test ./... -cover

# 特定パッケージ
go test ./internal/domain -v
go test ./internal/usecase -v
```

### 開発コマンド

```bash
# 開発サーバー起動（ホットリロード）
wails dev

# フロントエンド単体開発
cd frontend
npm run dev      # Next.js開発サーバー
npm run lint     # Biome lint
npm run format   # Biome format

# 本番ビルド
wails build

# Wails binding再生成（Go構造体変更時）
wails generate module

# Go依存関係更新
go mod tidy
```

### Git操作

```bash
# リリース作成（タグプッシュでGitHub Actions起動）
git tag v1.0.0
git push origin v1.0.0

# GitHub ActionsがDMGファイルをビルドしてリリースページに公開
```

## コード修正ガイドライン

### 1. ドメイン層の変更

**手順**:
1. `internal/domain/*_test.go`にテストを追加（Red）
2. 最小限の実装（Green）
3. リファクタリング（Refactor）

**例**: 新しいStrategy追加
```go
// 1. テストを先に書く
func TestNewStrategy(t *testing.T) {
    strategy := NewCustomStrategy(...)
    result := strategy.Apply("test.txt")
    assert.Equal(t, "expected.txt", result)
}

// 2. RenameStrategyを実装
type CustomStrategy struct {}
func (s *CustomStrategy) Apply(filename string) string {
    // 実装
}
```

### 2. ユースケース層の変更

**重要**: モックを使用してテスト

```go
type MockFileSystemService struct {
    mock.Mock
}

func TestRenameUseCase(t *testing.T) {
    mockFS := new(MockFileSystemService)
    mockFS.On("RenameFile", ...).Return(nil)

    useCase := NewRenameUseCase(mockFS)
    // テスト実行

    mockFS.AssertExpectations(t)
}
```

### 3. フロントエンド修正

**CSS変数の追加**:
```css
/* globals.css */
:root {
  --new-color: #...;
}
[data-theme='dark'] {
  --new-color: #...;
}
```

**コンポーネント修正後**:
- `dark:`クラスは使用しない（CSS変数を使う）
- `bg-background`, `text-foreground`等のユーティリティクラスを使用

### 4. Go構造体変更時

**重要**: Wails bindingの再生成が必要

```bash
# app.goの公開メソッドや構造体を変更した場合
wails generate module

# フロントエンドで新しい型が使える
```

### 5. 状態管理の注意点

**RenamePanel.tsxの重要な状態**:
- `currentFiles`: リネーム後も保持（連続リネーム対応）
- `message`: 自動クリアしない（ユーザーフィードバック重視）
- `previews`: ファイル選択時に即座に生成

## デバッグとトラブルシューティング

### フロントエンドデバッグ

```bash
# ブラウザで開発
wails dev
# http://localhost:34115 にアクセス
# Chrome DevToolsでデバッグ可能
```

### バックエンドデバッグ

```go
// ログ出力
import "log"
log.Printf("Debug: %v", value)

// app.goでエラーをフロントエンドに返す
if err != nil {
    return nil, fmt.Errorf("エラー詳細: %w", err)
}
```

### よくある問題

1. **プレビューが表示されない**
   - `GeneratePreview`のエラーログを確認
   - 正規表現が不正な可能性

2. **履歴が保存されない**
   - `~/.config/rename/`ディレクトリの書き込み権限確認
   - `JSONHistoryRepository`のエラーログ確認

3. **ビルドエラー**
   - `go mod tidy`実行
   - `cd frontend && npm ci`で依存関係再インストール

## CI/CD

### GitHub Actions

**トリガー**: タグプッシュ（`v*`）

**ワークフロー** (`.github/workflows/release.yml`):
1. Go/Node.jsセットアップ
2. フロントエンド依存関係インストール
3. Wailsビルド（Universal Binary: Intel + Apple Silicon）
4. DMG作成
5. GitHubリリース作成

**リリース手順**:
```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actionsが自動実行
# Releasesページに DMG がアップロード
```

## コーディング規約

### Go

- `gofmt`でフォーマット
- エラーは必ず処理
- 公開APIはコメント必須
- テストファイル名: `*_test.go`

### TypeScript

- Biomeでlint/format
- 型は明示的に定義
- `any`は避ける
- コンポーネントは関数コンポーネント

### 命名規則

- **Go**: PascalCase (公開), camelCase (非公開)
- **TypeScript**: PascalCase (コンポーネント/型), camelCase (変数/関数)
- **ファイル**: snake_case (Go), PascalCase (React)

## パフォーマンス最適化

### フロントエンド

- Debounce: 300ms（プレビュー生成）
- React.memo: 不要（小規模アプリ）
- 履歴: 最大10件表示（パフォーマンス維持）

### バックエンド

- ファイル操作: 並列処理不要（順次実行で十分）
- メモリ: ファイルパス文字列のみ保持

## セキュリティ

- ファイルパス: バリデーション不要（OS制限に依存）
- 履歴: ローカル保存のみ、外部送信なし
- 権限: ユーザーがアクセス可能なファイルのみ

## 将来の拡張ポイント

1. **新しいRenameStrategy追加**
   - `RenameStrategy`インターフェースを実装
   - テスト追加
   - UIにオプション追加

2. **プラグインシステム**
   - Strategy動的ロード
   - カスタムルール定義

3. **バッチ処理**
   - フォルダ一括処理
   - 再帰的リネーム

4. **履歴拡張**
   - タグ付け
   - お気に入り機能

---

このドキュメントは、AI開発者および人間の開発者がプロジェクトを理解し、一貫性のある変更を加えるためのガイドです。
