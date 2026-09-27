# CLAUDE.md

## 開発手法

- TDD（Test-Driven Development）を必須とする。Red（テストを先に書き、失敗することを確認）→ Green（最小限の実装で通す）→ Refactor の順で進める。
- テストツールは testify（assert, mock）。ユースケース層はインターフェースをモックしてテストする。

## レイヤーの依存規則

- ドメイン層（`internal/domain/`）は他のレイヤーに依存しない。
- プレゼンテーション層（`app.go`）はビジネスロジックを含まず、UseCase に処理を委譲する。
- UseCase はインターフェース（`FileSystemService`, `HistoryRepository`）に依存し、具象クラスは `NewApp()` で注入する。

## フロントエンド修正

- `dark:`クラスは使用しない（CSS変数を使う）
- `bg-background`, `text-foreground`等のユーティリティクラスを使用

## Go構造体変更時

**重要**: `app.go`の公開メソッドや構造体を変更した場合、Wails bindingの再生成が必要

```bash
wails generate module
```

## 状態管理の注意点

**RenamePanel.tsxの重要な状態**:
- `currentFiles`: リネーム後も保持（連続リネーム対応）
- `message`: 自動クリアしない（ユーザーフィードバック重視）
- `previews`: ファイル選択時に即座に生成

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
