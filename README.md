# Run
```bash
$ soi list
 - 01:  google          https://www.google.co.jp                      [search] 
 - 02:  yahoo           https://www.yahoo.co.jp                       [search] 
```

```bash
$ soi open yahoo
  -> (Chrome上で yahooが開きます)
```

```bash
$ soi add wikipedia https://ja.wikipedia.org/wiki/
```


# Install
## install golang
```bash
$ brew install go
```
## Build
```bash
$ go build -o soi
```

## Make config
```bash
$ mv sois.json.template sois.json
```

## Deploy
```bash
$ cp soi /usr/local/soi
``` 

TODO:
複数ページ立ち上げ
複数ページ管理
