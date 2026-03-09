APP_NAME ?= hikvision-backup
RELEASE_BRANCH ?= master

# Get current version from the latest git tag
VERSION ?= $(shell git describe --tags --abbrev=0)

# Strip the 'v' prefix if present
RAW_VERSION := $(subst v,,$(VERSION))

# Split version into parts (major.minor.patch)
MAJOR := $(word 1, $(subst ., ,$(RAW_VERSION)))
MINOR := $(word 2, $(subst ., ,$(RAW_VERSION)))
PATCH := $(word 3, $(subst ., ,$(RAW_VERSION)))

# Increment patch version by one
NEW_PATCH := $(shell echo $$(( $(PATCH) + 1 )))
NEW_VERSION_PATCH := v$(MAJOR).$(MINOR).$(NEW_PATCH)

# Increment minor version and reset patch to 0
NEW_MINOR := $(shell echo $$(( $(MINOR) + 1 )))
NEW_VERSION_MINOR := v$(MAJOR).$(NEW_MINOR).0

# Increment major version and reset minor/patch to 0
NEW_MAJOR := $(shell echo $$(( $(MAJOR) + 1 )))
NEW_VERSION_MAJOR := v$(NEW_MAJOR).0.0

NEW_VERSION_ALPHA := $(NEW_VERSION_PATCH)-alpha
NEW_VERSION_BETA := $(NEW_VERSION_PATCH)-beta

COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

I18N_TOOL_DIR ?= tools/readme-i18n-sync
I18N_SOURCE ?= ../../README.md
I18N_OUT_DIR ?= ../../i18n
I18N_TM_DIR ?=
I18N_PUBLISH_REMOTE ?= readme-i18n-sync
I18N_RUN_BASE = cd $(I18N_TOOL_DIR) && go run ./cmd/readme-i18n-sync --source $(I18N_SOURCE) --i18n-dir $(I18N_OUT_DIR) $(if $(I18N_TM_DIR),--tm-dir $(I18N_TM_DIR),)

LDFLAGS_STRING = -X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.buildDate=${DATE}'
LDFLAGS = -ldflags="${LDFLAGS_STRING}"

define CHECK_ON_RELEASE_BRANCH
@if [ "$(shell git rev-parse --abbrev-ref HEAD)" != "$(RELEASE_BRANCH)" ]; then \
	echo "Error: $(1) must be run on $(RELEASE_BRANCH) (current: $(shell git rev-parse --abbrev-ref HEAD))."; \
	exit 1; \
fi
endef

define CHECK_NOT_ON_RELEASE_BRANCH
@if [ "$(shell git rev-parse --abbrev-ref HEAD)" = "$(RELEASE_BRANCH)" ]; then \
	echo "Error: $(1) must not be run on $(RELEASE_BRANCH) (current: $(shell git rev-parse --abbrev-ref HEAD))."; \
	exit 1; \
fi
endef

define CHECK_NO_VERSION_TAG
@if git tag --points-at HEAD | grep -Eq '^v'; then \
	echo "Error: current commit already has a version tag (v...)."; \
	exit 1; \
fi
endef

define CHECK_TAG_NOT_EXISTS
@if git tag --list '$(1)' | grep -Eq '.'; then \
	echo "Error: tag $(1) already exists in repository."; \
	exit 1; \
fi
endef

define CHECK_TAG_EXISTS
@if git tag --list '$(1)' | grep -Eq '.'; then \
	:; \
else \
	echo "Error: tag $(1) does not exist locally."; \
	exit 1; \
fi
endef

define CHECK_TAG_ON_HEAD
@if git tag --points-at HEAD | grep -Fxq '$(1)'; then \
	:; \
else \
	echo "Error: tag $(1) is not on HEAD (current commit)."; \
	exit 1; \
fi
endef

define CHECK_PRE_RELEASE
	@$(MAKE) fmt-check
	@$(MAKE) i18n-check
endef

.PHONY: build test fmt fmt-check i18n-tool-test i18n-update i18n-check i18n-sync i18n-publish prepare-release next-alpha next-beta next-patch next-minor next-major release
# Build binary with version/commit/date baked via ldflags
build:
	@echo "==> Building ${APP_NAME}..."
	go build ${LDFLAGS} -o ${APP_NAME} .

# Run tests
test:
	@echo "==> Running tests..."
	go test ./...

# Format Go files
fmt:
	@echo "==> Formatting Go files..."
	gofmt -w .

# Check Go formatting
fmt-check:
	@echo "==> Checking Go formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "gofmt issues:"; \
		gofmt -l .; \
		exit 1; \
	fi

# Update translations (requires DEEPL_API_KEY)
i18n-update:
	@echo "==> Updating translations..."
	$(I18N_RUN_BASE) $(if $(I18N_FORCE),--force)

# Check translations (no API calls)
i18n-check:
	@echo "==> Checking translations..."
	$(I18N_RUN_BASE) --check

# Run tests for the standalone i18n module
i18n-tool-test:
	@echo "==> Testing readme-i18n-sync module..."
	cd $(I18N_TOOL_DIR) && go test ./...

# Update translations and commit changes (run before tagging)
i18n-sync:
	@echo "==> Updating translations and committing..."
	$(I18N_RUN_BASE) $(if $(I18N_FORCE),--force)
	@$(MAKE) i18n-check
	@if ! git diff --quiet; then \
		git add i18n/README.*.md i18n/tm/README.*.json; \
		git commit -m "docs: update translations"; \
	else \
		echo "No translation changes to commit."; \
	fi

# Publish nested module subtree to standalone repository
i18n-publish:
	@echo "==> Publishing readme-i18n-sync to $(I18N_PUBLISH_REMOTE)..."
	REPO_REMOTE=$(I18N_PUBLISH_REMOTE) $(I18N_TOOL_DIR)/scripts/publish-subtree.sh

.PHONY: bump
# Update deps + vendor + commit in one step
bump:
	@echo "==> Updating dependencies and vendor..."
	go get -u ./...
	go mod tidy
	go mod vendor
	git add go.mod go.sum vendor
	git commit -m "bump"

# Tag next patch version with alpha suffix
next-alpha:
	$(call CHECK_NOT_ON_RELEASE_BRANCH,next-alpha)
	$(call CHECK_NO_VERSION_TAG)
	$(call CHECK_TAG_NOT_EXISTS,$(NEW_VERSION_ALPHA))
	$(call CHECK_PRE_RELEASE)
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_ALPHA)..."
	git tag $(NEW_VERSION_ALPHA)

# Tag next patch version with beta suffix
next-beta:
	$(call CHECK_NOT_ON_RELEASE_BRANCH,next-beta)
	$(call CHECK_NO_VERSION_TAG)
	$(call CHECK_TAG_NOT_EXISTS,$(NEW_VERSION_BETA))
	$(call CHECK_PRE_RELEASE)
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_BETA)..."
	git tag $(NEW_VERSION_BETA)

# Tag next patch version
next-patch:
	$(call CHECK_ON_RELEASE_BRANCH,next-patch)
	$(call CHECK_NO_VERSION_TAG)
	$(call CHECK_TAG_NOT_EXISTS,$(NEW_VERSION_PATCH))
	$(call CHECK_PRE_RELEASE)
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_PATCH)..."
	git tag $(NEW_VERSION_PATCH)

# Tag next minor version (patch = 0)
next-minor:
	$(call CHECK_ON_RELEASE_BRANCH,next-minor)
	$(call CHECK_NO_VERSION_TAG)
	$(call CHECK_TAG_NOT_EXISTS,$(NEW_VERSION_MINOR))
	$(call CHECK_PRE_RELEASE)
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_MINOR)..."
	git tag $(NEW_VERSION_MINOR)

# Tag next major version (minor/patch = 0)
next-major:
	$(call CHECK_ON_RELEASE_BRANCH,next-major)
	$(call CHECK_NO_VERSION_TAG)
	$(call CHECK_TAG_NOT_EXISTS,$(NEW_VERSION_MAJOR))
	$(call CHECK_PRE_RELEASE)
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_MAJOR)..."
	git tag $(NEW_VERSION_MAJOR)

# Release: push release branch and current tag
release:
	$(call CHECK_ON_RELEASE_BRANCH,release)
	$(call CHECK_TAG_EXISTS,$(VERSION))
	$(call CHECK_TAG_ON_HEAD,$(VERSION))
	$(call CHECK_PRE_RELEASE)
	@echo "==> Releasing ${APP_NAME} version $(VERSION)..."
	git push origin $(RELEASE_BRANCH)
	git push origin $(VERSION)

# Prepare for release (run on release branch before tagging)
prepare-release:
	$(call CHECK_ON_RELEASE_BRANCH,prepare-release)
	$(call CHECK_NO_VERSION_TAG)
	@$(MAKE) fmt-check
	@$(MAKE) i18n-sync
