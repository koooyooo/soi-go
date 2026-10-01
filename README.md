<p align="center">
  <img src="./_img/soi-logo.jpg">
</p>

# soi
[![test](https://github.com/koooyooo/soi-go/actions/workflows/test.yaml/badge.svg)](https://github.com/koooyooo/soi-go/actions/workflows/test.yaml)
[![lint](https://github.com/koooyooo/soi-go/actions/workflows/lint.yaml/badge.svg)](https://github.com/koooyooo/soi-go/actions/workflows/lint.yaml)

**English** | [日本語](./README_JA.md)

## Overview
`soi` is a fast, local, CLI-based bookmark manager. You can add, search, and open bookmarks from a `soi>` prompt. Bookmarks are stored under `${HOME}/.soi`. Syncing that directory with cloud storage (or similar) also works across machines.

## Features
- **Add**: `soi> add {dir} {name} {url}`
- **Search**: `soi> list` to list and filter bookmarks
- **Open**: select from `list`, or use `open`
- **Browse by path**: `soi> dig` to walk directories
- **Buckets**: `soi> cb {name}` to switch contexts (e.g. work / hobby)

## Install

### Requirements
- Go 1.21 or later

### Setup

#### Config
Create `${HOME}/.soi/config.json` if needed. If it is missing, a default file is created on startup (mode `0600`).

```json
{
  "default_bucket": "default",
  "default_repository": "file",
  "default_browser": "firefox"
}
```

> `default_browser`: `firefox` | `chrome` | `safari` | `edge`  
> `default_repository`: `file` (default) or `sqlite`  
> Do not put the sync password in the config file. Use `SOI_USER_PASS` (and optionally `SOI_USER_NAME` / `SOI_SERVER`).

#### Install the binary
```
$ go install github.com/koooyooo/soi-go@latest
```

### Start
```bash
~$ soi-go
soi>
```

### `add`
Add a bookmark. `{dir}` and `{name}` are optional.

```bash
soi> add {dir} {name} https://www.google.com
```

- Words starting with `#` are tags (usable in `list` / `dig` filtering):

```bash
soi> add {dir} {name} https://www.google.com #search #entry
```

> #### Options
>
> - `-d`: directory (default: `YYYY-MM`, e.g. `2026-10`)
> - Directories may be nested with `/`
> ```
> soi> add -d search https://www.google.com
> ```
>
> - `-n`: name (default: the page `<title>`)
> ```
> soi> add -n google https://www.google.com
> ```

### `list`
List bookmarks, filter them, then open one.

- Use `Tab` / `Shift+Tab` or `↑` / `↓` to move selection:

```bash
soi> list
           adf46ead [ 10 00.0] api/API設計ガイド [#guide #api]
           a73764ed [  1 00.0] books/GooglePlay-Audiobooks
           46734e6c [  6 00.0] contents/MDN [#guide]
```

- Type part of a path, name, or tag to filter.

#### Filter by name
```bash
soi> list MDN
           46734e6c [  6 00.0] contents/MDN [#guide]
```

#### Filter by tag
```bash
soi> list #guide
           adf46ead [ 10 00.0] api/API設計ガイド [#guide #api]
           46734e6c [  6 00.0] contents/MDN [#guide]
```

- Press `Enter` on a row to open it in a browser.

> Browser options (otherwise the config default is used):
> - `-c` chrome
> - `-f` firefox
> - `-s` safari
> - `-e` edge

> Sort options:
> - `-n` view count
> - `-a` added date (newest first)
> - `-v` last viewed (newest first)

### `dig`
Walk directories hierarchically. Prefer `list` for global search; use `dig` to drill down.

```bash
soi> dig
          search/
          sns/
```

- `Tab` / `↓` to select a directory
- `→` to list its contents

```bash
soi> dig search/
                 search/google.json
                 search/yahoo.json
```

- `Enter` opens the bookmark in a browser

### `mv` / `rm`
Move or remove within the current bucket. Paths outside the bucket are rejected.  
You can use the same logical path as in `list` (no extension), or the on-disk `.json` name.  
If deletion fails, check `(bucket: ...)` in the error and switch with `cb` if needed.

```bash
soi> mv search/google archive/google
soi> rm archive/google
# or
soi> rm archive/google.json
```

### `cb`
Switch buckets (contexts such as `work` or `hobby`).

```bash
soi> cb hobby
create & change current bucket: hobby
```

```bash
soi> cb hobby
change current bucket: hobby
```

```bash
soi> cb
current bucket: [hobby]
```

### `quit`
Leave the `soi>` prompt.

```
soi> quit
$
```
