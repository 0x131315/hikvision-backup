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

## Make Targets

```bash
make test
make vendor
make vendor-check
make check SOURCE=../../README.md I18N_DIR=../../i18n
make update SOURCE=../../README.md I18N_DIR=../../i18n
make sync SOURCE=../../README.md I18N_DIR=../../i18n
make print-next-tag
make tag TAG=readme-i18n-sync/v0.1.0
make release-tag TAG=readme-i18n-sync/v0.1.0 REMOTE=readme-i18n-sync
```

If `TAG` is not set, module make targets use automatic patch increment from the latest `readme-i18n-sync/v*` tag.
The module uses vendored dependencies by default (`GOFLAGS=-mod=vendor`).
