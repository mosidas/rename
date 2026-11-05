# rename

macOS で、選んだ複数のファイルの名前を文字列の置換でまとめて変えるアプリ。

## 動作環境

- macOS(Apple Silicon と Intel の両方)

## インストール

[Releases](https://github.com/mosidas/rename/releases) から `rename-<バージョン>.zip` をダウンロードして展開し、`rename.app` を `/Applications` へ移して開く。Apple Silicon と Intel の両方で動くユニバーサルビルドである。

配布物はアドホック署名だけで公証していないため、初回の起動は Gatekeeper に止められる。次のどちらかで起動を許可する。

- ターミナルで `xattr -dr com.apple.quarantine /Applications/rename.app` を実行し、隔離属性を外す。
- 一度開いてブロックされたあと、システム設定 > プライバシーとセキュリティ の下部にある「このまま開く」を押す(macOS 14 では Finder で右クリックして「開く」でもよい)。

ソースからビルドするときは、Go 1.26 以上・Node.js 24 以上・[Wails CLI](https://wails.io/docs/gettingstarted/installation) v2 を入れたうえで `scripts/bundle.sh` を実行し、できた `build/bin/rename.app` を `/Applications` へコピーして開く。

どちらの方法でも、開くと「ファイルを選択」のボタンを持つウィンドウが出る。

## 使い方

1. 「ファイルを選択」を押し、名前を変えるファイルを複数選ぶ。ファイルを引数にして起動しても、同じように読み込まれる。
2. 「置換前」に探す文字列、「置換後」に置き換える文字列を入れる。入れるたびに、元のファイル名と新しいファイル名の一覧が更新される。
3. 必要なら「正規表現」「大文字小文字を区別しない」にチェックを付ける。
4. 「リネーム実行」を押す。成功と失敗の件数が表示される。

使った置換のパターンは履歴として `~/.config/rename/config.json` に保存され、次回から「置換前」の欄を選択すると一覧から選び直せる。履歴を削除するときはこのファイルを削除する。

Finder で選んだファイルを右クリックのメニューから開く手順は [docs/finder-quick-action.md](docs/finder-quick-action.md) にある。

## 開発

Go 1.26 以上・Node.js 24 以上・[Wails CLI](https://wails.io/docs/gettingstarted/installation) v2 を入れる。

```sh
scripts/bundle.sh   # build/bin/rename.app を作る
go test ./...
```

## ライセンス

[MIT License](LICENSE)
