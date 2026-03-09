#!/usr/bin/env bash
set -euo pipefail

REPO_REMOTE="${REPO_REMOTE:-readme-i18n-sync}"
PREFIX="tools/readme-i18n-sync"
BRANCH="readme-i18n-sync-release"

if ! git rev-parse --git-dir >/dev/null 2>&1; then
  echo "Run from a git repository" >&2
  exit 1
fi

git subtree split --prefix="${PREFIX}" -b "${BRANCH}"
git push "${REPO_REMOTE}" "${BRANCH}:main"
git branch -D "${BRANCH}"

echo "Published ${PREFIX} to ${REPO_REMOTE}:main"
