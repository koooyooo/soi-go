<p align="center">
  <img src="./_img/soi-logo.jpg">
</p>

# soi
[![test](https://github.com/koooyooo/soi-go/actions/workflows/test.yaml/badge.svg)](https://github.com/koooyooo/soi-go/actions/workflows/test.yaml)
[![lint](https://github.com/koooyooo/soi-go/actions/workflows/lint.yaml/badge.svg)](https://github.com/koooyooo/soi-go/actions/workflows/lint.yaml)

[English](./README.md) | **日本語**


## 概要
`soi` は CLIベースの ローカルで高速なブックマークマネージャです。機能としてはブックマークの追加・検索・閲覧が存在します。追加されたブックマークはローカルストレージ（`${HOME}/.soi`）に管理されますが、このストレージをクラウドストレージと連携することで、異なるPC間の同期も可能です。

## 主な機能
- **ブックマークの追加**: `soi> add {dir} {name} {url}` でブックマークを追加できます。
- **ブックマークの検索**: `soi> list` でブックマークをリストアップし、任意のキーワードで絞り込めます。
- **ブックマークの閲覧**: `soi> open {id}` / `list` で選択して閲覧できます。
- **階層探索**: `soi> dig` でディレクトリを掘り下げます。
- **バケット切替**: `soi> cb {name}` でコンテキストを切り替えます。


## インストール

### 必要要件
- `go` がインストールされていること（1.21 以上）

### インストール手順

#### 設定ファイル
`${HOME}/.soi/config.json` を作成します。存在しない場合は起動時に既定ファイルが生成されます（権限 `0600`）。

```json
{
  "default_bucket": "default",
  "default_repository": "file",
  "default_browser": "firefox"
}
```
> `default_browser` の選択肢は `firefox` | `chrome` | `safari` | `edge` です  
> `default_repository` は `file`（既定）または `sqlite` です  
> サーバ同期用のパスワードは設定ファイルに書かず、環境変数 `SOI_USER_PASS` を使います（`SOI_USER_NAME` / `SOI_SERVER` も可）

#### バイナリインストール 
`$ go install` でインストールします
```
$ go install github.com/koooyooo/soi-go@latest
```

ローカルで開発する場合は `make install` でも同じ場所（`$(go env GOPATH)/bin`）に `soi-go` を配置します。

### 起動
`$ soi-go` と打ち込むと `soi>` 形式のプロンプトが立ち上がります

```bash
~$ soi-go
soi>
```

### `add`
`add`(追加) コマンドでブックマークを追加します
- ブックマークには分類用のディレクトリと識別用の名前を指定可能です。それぞれ`{dir}` `{name}` で指定します。これらは省略可能です。
```bash
soi> add　{dir} {name} https://www.google.com
```

- `#`で開始された用語はタグとなります
  - タグは後述の `list`, `dig` コマンドにおける絞り込みで活用可能です。
```bash
soi> add　{dir} {name} https://www.google.com #search #entry
```

> #### options
> オプションで各要素を明示的に指定できます。
> 
> - `-d`オプションでディレクトリを明示的に指定できます
> - 省略時のデフォルトディレクトリは `YYYY-MM`（例: `2026-10`）です
> - ディレクトリは `/`区切りで階層的に表現することも可能です
> ```
> soi> add -d search https://www.google.com
> ```
> 
> - `-n`オプションで名前を明示的に指定できます
> - `-n`オプションが無い場合のデフォルト値は 対象URL内の `<title>`要素です
> ```
> soi> add -n google https://www.google.com
> ```


### `list`
`list` コマンドで ブックマークをリストアップし、絞り込み、最後に選択することが可能です。

- `Tab`（`Shift + Tab`）キーや `↑`・`↓`キーで対象を選択します
```bash
soi> list 
           adf46ead [ 10 00.0] api/API設計ガイド [#guide #api]                               
           a73764ed [  1 00.0] books/GooglePlay-Audiobooks                      
           46734e6c [  6 00.0] contents/MDN [#guide]                                    
```

- ディレクトリや名前の一部, タグ名をタイプすることで対象の絞り込みを行います

#### 名前で絞り込み
```bash
soi> list MDN
           46734e6c [  6 00.0] contents/MDN [#guide]                                    
```

#### タグで絞り込み
```bash
soi> list #guide
           adf46ead [ 10 00.0] api/API設計ガイド [#guide #api]                               
           46734e6c [  6 00.0] contents/MDN [#guide]                                    
```

- 選択した行で`Enter`を押下すると、対象のブックマークをブラウザで開きます

> Note: リスト時に**ブラウザ指定**オプションを付けることで指定したブラウザが開きます。それ以外の場合は設定ファイルで指定したデフォルトブラウザが起動します。
> - `-c` `chrome`
> - `-f` `firefox`
> - `-s` `safari`
> - `-e` `edge`

> Note: リスト時に**ソート指定**オプションをつけることで指定した順にソートされます。
> - `-n` 閲覧回数
> - `-a` 追加日時（新しい順）
> - `-v` 閲覧日時（新しい順）

### `dig`
`dig` コマンドで、ディレクトリを階層的に探索します。
`list` が全体検索を行うイメージなのに対し、digは階層（ディレクトリ）を掘り下げてゆくイメージです。

```bash
soi> dig
          search/
          sns/
```
- `Tab` キーや `↓` キーでディレクトリを選択します
- 選択した状態で `→` を選ぶと内部要素が提示されます

```bash
soi> dig search/
                 search/google.json
                 search/yahoo.json
``` 
- `Enter` キーを押下すると、ブラウザでブックマークを開きます

### `mv` / `rm`
同一バケット内での移動・削除です。バケット外へのパス脱出は拒否されます。
`list` と同じ論理パス（拡張子なし）でも指定できます。実ファイル名（`.json`）でも構いません。
削除できないときは、エラー末尾の `(bucket: ...)` で現在バケットを確認してください（`cb` で切替）。

```bash
soi> mv search/google archive/google
soi> rm archive/google
# または
soi> rm archive/google.json
```

### `cb`
`cb` コマンドで バケットを切り替えます。バケットとはブックマークをコンテキスト毎に整理するためのもので、例えば `work`, `hobby` 等です。
- `cb` コマンドに続けてバケット名を入力します
- 初回は既存のバケットが無いため新規に作成した上で切り替えます
```bash
soi> cb hobby
create & change current bucket: hobby
```
- 既存のものがあれば、単純に切り替えます
```bash
soi> cb hobby
change current bucket: hobby
```
- 引数なしで実行すると、現在のバケットを確認できます
```bash
soi> cb
current bucket: [hobby]
```
### `quit`
`quit` コマンドで `soi`プロンプトを抜けることができます

```
soi> quit
$
```
