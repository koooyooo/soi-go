
# Soi
Soi CLIベースのURL管理するツールです。 CLIベースなので操作性に慣れれば慣れるほど快適でシンプルな操作を行うことができます。  
また、データはローカルにストアされるためレスポンスはクイックです。

## 導入
任意の ディレクトリに [git](https://git-scm.com/) ベースでクローンします。[make](https://www.gnu.org/software/make/) 経由でインストールを行います。
```
$ git clone https://github.com/koooyooo/soi-go.git
$ cd soi-go && make install
```

## 利用法
### 立ち上げ
`soi`と打ち込むとSoiが立ち上がり、`soi>`形式のプロンプトを立ち上げてユーザにコマンド入力を促します。soiの操作を楽しみましょう。

```
$ soi
soi> 
```

#### 追加
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

#### リスト
Type `list` to list up the links(sois).  
- Then, you can choose a link by using `Tab` key.
- Next, you can select links by using `Tab`, `Shift + Tab` or arrow keys.
- Also, you can filter links by typing free words.
- Finally, type `enter` key to open the site in a browser.
- Default browser is `chrome`, if you like `firefox`, use `-f` option, or you prefer `safari`, use `-s` option.
```
soi> list
           [  0 00.0] search/google.json
           [  0 00.0] search/yahoo.json
           [  0 00.0] sns/facebook.json
```


#### Dig
Type `Dig` to search links(sois) by digging directory one by one.
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
