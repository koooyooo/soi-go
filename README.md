
![soi](./soi.jpg)

[![basic](https://github.com/koooyooo/soi-go/actions/workflows/basic.yaml/badge.svg)](https://github.com/koooyooo/soi-go/actions/workflows/basic.yaml)

# soi
`soi` は `golang`製の CLIによるブックマークマネージャです。CLIによる快適な操作が可能です。データはローカルストレージに保存されます。ローカルストレージとしては `${HOME}/.soi` が割り当てられます。

## Install
`$ go install` 
```
$ go install github.com/koooyooo/soi-go@latest
```

## Usage
### `soi` >
`soi`と打ち込むと `soi>`形式のプロンプトが立ち上がります。
```
$ soi
soi> 
```

### `add`
`add`(追加) コマンドでブックマークを追加します
- ブックマークには分類用のディレクトリと識別用の名前を指定可能です。それぞれ`{dir}` `{name}` で指定します。これらは省略可能です。
```
soi> add　{dir} {name} https://www.google.com
```

- `-d`オプションでディレクトリを明示的に指定できます
- 省略時のデフォルトディレクトリは `new`です
- ディレクトリは `/`区切りで階層的に表現することも可能です
```
soi> add -d search https://www.google.com
```

- `-n`オプションで名前を明示的に指定できます
- `-n`オプションが無い場合のデフォルト値は 対象URL内の `<title>`要素です
```
soi> add -n google https://www.google.com
```


### `list`
`list` コマンドで ブックマークをリストアップすることが可能です

- `Tab`（`Shift + Tab`）キーや `↑`・`↓`キーで対象を選択します
- ディレクトリや名前の一部をタイプすることで対象の絞り込めます
- 選択した行で`Enter`を押下すると、ブラウザでブックマークを開きます
```
soi> list 
           adf46ead [ 10 00.0] api/API設計ガイド                                
           a73764ed [  1 00.0] books/GooglePlay-Audiobooks                      
           46734e6c [  6 00.0] contents/MDN                                     
```

### dig（掘り下げ）
`dig`（掘り下げ）コマンドで、ディレクトリを階層的に探索することが可能です。`list`と同様に最終的には URLを開きますが、listが全体検索を行うイメージなのに対し、digはジャンルごとにコンテンツを掘り下げてゆくイメージです。

```
soi> dig
          search/
          sns/
```
- Then, type `Tab` key and choose one directory
- Next, type `→` key to dig the chosen directory
```
soi> dig search/
                 search/google.json
                 search/yahoo.json
``` 
- Finally, type `enter` key to open the site in a browser
- Default browser is `chrome`, if you like `firefox`, use `-f` option, or you prefer `safari`, use `-s` option.

<!--
### tag              
```
  Desc:  not implemented now
```

### mv
```
  Desc:  move file to dir 
  Usage: mv (current path) to (dir)
```
-->

### `quit`
`quit` コマンドで `soi`プロンプトを抜けることができます

```
soi> quit
$
```
