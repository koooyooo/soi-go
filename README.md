

## Usage

### Add
`Add` add a link
```bash
# No Tags
$ soi add wikipedia https://ja.wikipedia.org/wiki/

# With Tags
$ soi add -t "web,lib" wikipedia https://ja.wikipedia.org/wiki/
```

### List
`list` list up links
```bash
$ soi list
 - 01:  google          https://www.google.co.jp                      [search cloud web]
 - 02:  yahoo           https://www.yahoo.co.jp                       [search]
 - 03:  youtube         https://www.youtube.com/?gl=JP                [video]
```

### Open
`open` open a link by google chrome
```bash
# open yahoo (by Google Chrome)
$ soi open yahoo
```

### Tag
`tag` add tags to a registered link
```bash
$ 
```
### Tags

### Remove


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
