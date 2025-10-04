![soi](./_img/soi-logo.png)

# Soi
Soi is a url management CLI system. which could add url, find url and rename url.
(Soi はURLを管理する CLIシステムです)

## Install
Clone this repository, then run `make install` in the directory.

```
$ git clone https://github.com/koooyooo/soi-go.git
$ cd soi-go && make install
```

## Usage

### Run Soi
Type `soi` (or `soi-go`: install by go install)  then soi prompt appears.
```
$ soi
soi> 
```

#### Add
Type `add` to add a new link (soi).  
```
soi> add https://www.google.com
```

- You can specify its directory by `-d` option (default: `new`)
```
soi> add -d search https://www.google.com
```

- You can specify its name by `-n` option (default: HTML's `<title>` value)
```
soi> add -n google https://www.google.com
```

- Specifying both `-d` and `-n` is the most preferable.
```
soi> add -d search -n google https://www.google.com
```

#### List
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
