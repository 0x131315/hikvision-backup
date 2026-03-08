# Release Process

This project uses git tags and GoReleaser. Release commands are defined in `Makefile`.

## Preconditions

- Ensure you are on the correct branch.
  - Release tags (`next-patch`, `next-minor`, `next-major`) must run on `RELEASE_BRANCH` (default: `master`).
  - Pre-release tags (`next-alpha`, `next-beta`) must NOT run on `RELEASE_BRANCH`.
- Ensure the current commit does not already have a `v*` tag.
- Ensure the target tag does not already exist in the repository.

## Tagging

Choose exactly one of the following:

- Patch release: `make next-patch`
- Minor release: `make next-minor`
- Major release: `make next-major`
- Alpha pre-release: `make next-alpha`
- Beta pre-release: `make next-beta`

Each command creates a git tag locally.

## Release (push)

Push the release branch and the current tag:

- `make release`

This runs:

- `git push origin $(RELEASE_BRANCH)`
- `git push origin $(VERSION)`

## Verification

- Confirm the tag exists locally: `git tag --list 'v*' | tail -n 5`
- Confirm the tag is on the expected commit: `git tag --points-at HEAD`
- Confirm the GitHub release workflow ran on the new tag.

## Notes

- The current version is derived from the latest git tag via `git describe --tags --abbrev=0`.
- `release` does not build binaries locally; it only pushes the branch and tag.
- If a tag was created on the wrong branch, delete it locally and remotely, then retag correctly.

## Translation Automation

Translations are managed via `tools/i18n` and updated automatically in CI.

Environment variables:
- **`DEEPL_API_KEY`** — DeepL API key (preferred).
- **`DEEPL_API_URL`** — Optional override for DeepL endpoint.
- **`GOOGLE_TRANSLATE_API_KEY`** — Google Cloud Translation API key (fallback).
- **`GOOGLE_TRANSLATE_API_URL`** — Optional override for Google Translate endpoint.
- **`LIBRETRANSLATE_URL`** — LibreTranslate base URL (fallback when `DEEPL_API_KEY` is not set).
- **`LIBRETRANSLATE_API_KEY`** — Optional API key for LibreTranslate instances that require it.
