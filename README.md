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

## Finder統合（Quick Action）

Finderから選択したファイルを右クリックメニューで直接Renameアプリで開くことができます。

### Automator Quick Actionの設定

#### 1. Automatorを起動

1. アプリケーション > Automator を開く
2. 「新規書類」→「Quick Action」を選択

#### 2. ワークフローを設定

1. **「ワークフローが受け取る現在の項目:」** を **「ファイルまたはフォルダ」** に設定
2. **「検索対象:」** を **「Finder.app」** に設定

#### 3. シェルスクリプトアクションを追加

1. 左側のアクションリストから「Run Shell Script」（シェルスクリプトを実行）を検索してダブルクリック
2. **「Pass input:」**（入力の引き渡し方法）を **「as arguments」**（引数として）に変更
3. 以下のスクリプトを入力:

```bash
for f in "$@"
do
    open -a "/Applications/rename.app" "$f"
done
```

#### 4. 保存

1. `⌘S` で保存
2. 名前: 「Renameで開く」（任意の名前でOK）
3. 保存先は自動的に `~/Library/Services/` になります

#### 5. 使用方法

1. Finderでファイルを選択（複数選択可）
2. 右クリック → 「Quick Actions」→ 「Renameで開く」を選択
3. Renameアプリが起動し、選択したファイルが自動的にロードされます

**ヒント**: 既にRenameアプリが起動している場合でも、新しいファイルが自動的にロードされ、既存のウィンドウが前面に表示されます。

### トラブルシューティング

**Quick Actionが表示されない場合**:
- システム設定 > プライバシーとセキュリティ > 機能拡張 > Finder拡張機能 で、Automatorが有効になっているか確認してください

**アプリのパスが違う場合**:
- アプリを `/Applications` 以外の場所にインストールした場合は、スクリプト内のパスを変更してください


## 技術スタック

- **Wails v2**
- **Go**
- **Next.js 15 + React 19**
- **TypeScript**
- **Tailwind CSS**

## ライセンス

[MIT License ](LICENSE)
