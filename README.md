
# Soi
Soi CLIベースのURL管理するツールです。 CLIベースなので操作性に慣れれば慣れるほど快適でシンプルな操作を行うことができます。
また、データはローカルに保管されるため、データ処理のレスポンスが早いのも特徴です。。

## 導入
任意の ディレクトリに [git](https://git-scm.com/) ベースでクローンします。[make](https://www.gnu.org/software/make/) 経由でインストールを行います。
```
$ git clone https://github.com/koooyooo/soi-go.git
$ cd soi-go && make install
```

## 利用法
### 立ち上げ
`soi`と打ち込むとSoiが立ち上がり、`soi>`形式のプロンプトを立ち上げてユーザにコマンド入力を促します。soiの操作を楽しんでください。

```
$ soi
soi> 
```

#### add (追加)
`add`(追加) コマンドでURLを追加することが可能です。引数に目的のURLを渡しましょう。
```
soi> add https://www.google.com
```

- URLは任意のディレクトリに分類して管理することができます。ディレクトリを指定する際には `-d`オプションを付けてください。`-d`オプションがない場合は `new`ディレクトリが既定となります。ディレクトリは `/`区切りで階層的に表現することも可能です。
```
soi> add -d search https://www.google.com
```

- URLには任意の名前をつけることが可能です。名前を指定する際には `-n`オプションを付けてください。`-n`オプションが無い場合には 対象のURLにアクセスし、HTML内の `<title>`属性値を取得して名前として利用します。
```
soi> add -n google https://www.google.com
```

- `-d`オプションと`-n`オプションは同時に付与することができます。これは最も望ましい利用法です。`-d`で登録時に細かやかな整理ができますし、`-n`で自分が使いやすい名前を指定できるからです。
```
soi> add -d search -n google https://www.google.com
```

#### list (リストアップ）
`list` (リストアップ) コマンドで URLをリストアップすることが可能です

- `Tab`キーで候補を選択することができます（`Shift + Tab`で逆順に選択します）
- `↑`・`↓`キーでも同じ操作が可能です。
- 任意の文字列を打ち込むことで、パスや名前で絞り込みをかけることが可能です
- 最後に`Enter`キーを打ち込むことで、ブラウザで当該 URLを開きます
```
soi> list 
           adf46ead [ 10 00.0] api/API設計ガイド                                
           a73764ed [  1 00.0] books/GooglePlay-Audiobooks                      
           46734e6c [  6 00.0] contents/MDN                                     
```

#### Dig（掘り下げ）
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


#### Tag              
```
  Desc:  not implemented now
```

#### Mv
```
  Desc:  move file to dir 
  Usage: mv (current path) to (dir)
```

#### Quit
```
  Desc:  quit soi> and go back to console $
  Usage: quit
```
