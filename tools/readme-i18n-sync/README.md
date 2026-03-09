# readme-i18n-sync

CLI tool to sync translated README files using translation memory blocks.

## Usage

From this repository root:

```bash
cd tools/readme-i18n-sync
go run ./cmd/readme-i18n-sync --source ../../README.md --i18n-dir ../../i18n --check
go run ./cmd/readme-i18n-sync --source ../../README.md --i18n-dir ../../i18n
go run ./cmd/readme-i18n-sync --source ../../README.md --i18n-dir ../../i18n --force
```

## Flags

- `--check` validate that all source blocks are translated
- `--init` bootstrap missing translations from source text
- `--force` retranslate all blocks and overwrite TM entries
- `--source` source README path (default `README.md`)
- `--i18n-dir` translated files directory (default `i18n`)
- `--tm-dir` translation-memory directory (default `<i18n-dir>/tm`)

## Publish To Standalone Repo

```bash
tools/readme-i18n-sync/scripts/publish-subtree.sh
```
