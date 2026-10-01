# feat-1: 保存経路の集約と欠陥修正

## 概要
ブックマークの追加・削除・移動を Repository / Service 経由に揃え、既知の欠陥（Path 補完、push、設定権限など）を直す。

## 要件
- `add` / `rm` / `mv` がシェル直叩きや生 `WriteFile` に依存しない
- `rm` はバケット配下のみ削除できる
- SQLite の `Remove` を実装する
- `LoadSoiData` がファイル位置から `Path` を正しく復元する
- `push` が HTTP 失敗時に落ちない
- `config.json` は `0600`、パスワードは環境変数優先
- execute 系の `flag.PanicOnError` と loader / completer の `log.Fatal` を error / 空候補に寄せる
- README / Makefile / テストを実装に合わせる

## 設計
- 書き込みの正は `service.Service` → `repository.Repository`（`add` は `Store`）
- `rm` / `mv` はバケット内パス操作だが `pathutil.ResolveUnderRoot` で脱出を防ぐ
- ファイル実体のパスは `SoiData.FilePath`（`ToStorableName`）+ repository `basePath`
- 読み込み時の `Path` は bucket root からの相対ディレクトリで正規化

## 実装結果
- ブランチ: `fix/repository-write-path`
- テスト: `go test ./...` 通過
