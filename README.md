# rename

macOS で、選んだ複数のファイルの名前を文字列の置換でまとめて変えるアプリ。

## 動作環境

- macOS(Apple Silicon と Intel の両方)

## インストール

### Releases から

1. [Releases](https://github.com/mosidas/rename/releases) から `rename-<バージョン>.zip` をダウンロードして展開する。
2. `rename.app` を `/Applications` へ移す。
3. 配布物はアドホック署名だけで公証していないため、初回の起動は Gatekeeper に止められる。ターミナルで `xattr -dr com.apple.quarantine /Applications/rename.app` を実行し、隔離属性を外す。
4. `rename.app` を開く。

### ソースからビルドする

1. Go 1.26 以上・Node.js 24 以上・[Wails CLI](https://wails.io/docs/gettingstarted/installation) v2 をインストールする。
2. `scripts/bundle.sh` を実行する。`build/bin/rename.app` ができる。
3. `rename.app` を `/Applications` へコピーして開く。

## 使い方

1. 「ファイルを選択」を押し、名前を変えるファイルを複数選択する。ファイルを引数にして起動しても、同じように読み込まれる。
2. 「置換前」に探す文字列、「置換後」に置き換える文字列を入力する。入力するたびに、元のファイル名と新しいファイル名の一覧が更新される。
3. 必要なら「正規表現」「大文字小文字を区別しない」にチェックを付ける。
4. 「リネーム実行」を押す。成功と失敗の件数が表示される。

使った置換のパターンは履歴として `~/.config/rename/config.json` に保存され、次回から「置換前」の欄を選択すると一覧から選び直せる。履歴を削除するときはこのファイルを削除する。

Finder の右クリックのメニューから開くには、[Finder で選んだファイルを rename で開く](docs/finder-quick-action.md) の手順でクイックアクションを作る。

## 開発

道具とビルドは「ソースからビルドする」の手順と同じである。テストは次のコマンドで実行する。

```sh
go test ./...
```

## ライセンス

[MIT License](LICENSE)
