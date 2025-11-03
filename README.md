# Rename - ファイル一括リネームツール

macOS用ファイル一括リネームアプリケーション。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Platform](https://img.shields.io/badge/platform-macOS-lightgrey.svg)


## インストール

### リリースからダウンロード

1. [Releases](https://github.com/mosidas/rename/releases)ページから最新のDMGファイルをダウンロード
2. DMGファイルをダブルクリック
3. `rename.app`を`Applications`フォルダにドラッグ＆ドロップ
4. アプリケーションフォルダから起動

### ソースからビルド

```bash
# 前提条件
# - Go 1.21以上
# - Node.js 20以上
# - Wails CLI v2

# リポジトリをクローン
git clone https://github.com/mosidas/rename.git
cd rename

# ビルド
wails build

# アプリケーションが build/bin/rename.app に生成されます
```

## 使い方

### 基本的な使い方

1. **ファイルを選択**
   - 「ファイルを選択」ボタンをクリック
   - リネームしたいファイルを複数選択

2. **置換パターンを入力**
   - 「置換前」: 検索したい文字列を入力
   - 「置換後」: 置き換え後の文字列を入力
   - プレビューがリアルタイムで更新されます

3. **オプション設定（必要に応じて）**
   - ✓ 正規表現: 正規表現パターンを使用
   - ✓ 大文字小文字を区別しない: 大小文字を無視

4. **実行**
   - 「リネーム実行」ボタンをクリック
   - 変更されたファイル数が表示されます


## 設定ファイル

履歴データは以下の場所に自動保存されます：

```
~/.config/rename/config.json
```

履歴をクリアしたい場合は、このファイルを削除してください。


## 技術スタック

- **Wails v2**
- **Go**
- **Next.js 15 + React 19**
- **TypeScript**
- **Tailwind CSS**

## ライセンス

[MIT License ](LICENSE)
