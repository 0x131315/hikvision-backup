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
I18N_SUBTREE_PREFIX ?= tools/readme-i18n-sync
I18N_SUBTREE_BRANCH ?= readme-i18n-sync-release
I18N_MODULE_TAG ?=
RELEASE_I18N_MODULE ?= 1
GOFMT_PATHS = -path ./vendor -o -path ./.git -o -path ./.cache -o -path ./bin

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

define CHECK_RELEASE_READY
	$(call CHECK_ON_RELEASE_BRANCH,release)
	$(call CHECK_TAG_EXISTS,$(VERSION))
	$(call CHECK_TAG_ON_HEAD,$(VERSION))
	$(call CHECK_PRE_RELEASE)
endef

.PHONY: build test fmt fmt-check i18n-tool-test i18n-update i18n-check i18n-sync i18n-subtree-split i18n-publish i18n-module-release release-preflight release-i18n-module release-push prepare-release next-alpha next-beta next-patch next-minor next-major release
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
	@find . \( $(GOFMT_PATHS) \) -prune -o -type f -name '*.go' -print0 | xargs -0r gofmt -w

# Check Go formatting
fmt-check:
	@echo "==> Checking Go formatting..."
	@fmt_issues="$$(find . \( $(GOFMT_PATHS) \) -prune -o -type f -name '*.go' -print0 | xargs -0r gofmt -l)"; \
	if [ -n "$$fmt_issues" ]; then \
		echo "gofmt issues:"; \
		echo "$$fmt_issues"; \
		exit 1; \
	fi

# Update translations (requires DEEPL_API_KEY)
i18n-update:
	@echo "==> Updating translations..."
	$(MAKE) -C $(I18N_TOOL_DIR) update SOURCE="$(I18N_SOURCE)" I18N_DIR="$(I18N_OUT_DIR)" TM_DIR="$(I18N_TM_DIR)" I18N_FORCE="$(I18N_FORCE)"

# Check translations (no API calls)
i18n-check:
	@echo "==> Checking translations..."
	$(MAKE) -C $(I18N_TOOL_DIR) check SOURCE="$(I18N_SOURCE)" I18N_DIR="$(I18N_OUT_DIR)" TM_DIR="$(I18N_TM_DIR)"

# Run tests for the standalone i18n module
i18n-tool-test:
	@echo "==> Testing readme-i18n-sync module..."
	$(MAKE) -C $(I18N_TOOL_DIR) test

# Update translations and commit changes (run before tagging)
i18n-sync:
	@echo "==> Updating translations and committing..."
	$(MAKE) -C $(I18N_TOOL_DIR) sync SOURCE="$(I18N_SOURCE)" I18N_DIR="$(I18N_OUT_DIR)" TM_DIR="$(I18N_TM_DIR)" I18N_FORCE="$(I18N_FORCE)"
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

# Create/update local subtree split branch for the i18n module
i18n-subtree-split:
	@echo "==> Splitting subtree $(I18N_SUBTREE_PREFIX) into branch $(I18N_SUBTREE_BRANCH)..."
	git subtree split --prefix=$(I18N_SUBTREE_PREFIX) -b $(I18N_SUBTREE_BRANCH)

# Publish module and push explicit module tag to standalone repository
# Usage: make i18n-module-release I18N_MODULE_TAG=readme-i18n-sync/v0.1.0
i18n-module-release:
	@split_sha="$$(git subtree split --prefix=$(I18N_SUBTREE_PREFIX))"; \
	remote_sha="$$(git ls-remote $(I18N_PUBLISH_REMOTE) refs/heads/main | awk '{print $$1}')"; \
	tag="$(I18N_MODULE_TAG)"; \
	if [ -z "$$tag" ]; then \
		tag="$$( $(MAKE) --no-print-directory -s -C $(I18N_TOOL_DIR) print-next-tag )"; \
	fi; \
	remote_tag_sha="$$(git ls-remote $(I18N_PUBLISH_REMOTE) refs/tags/$$tag | awk '{print $$1}')"; \
	if [ -n "$$remote_sha" ] && [ "$$split_sha" = "$$remote_sha" ]; then \
		if [ -n "$$remote_tag_sha" ]; then \
			echo "No module changes to release and tag $$tag already exists on $(I18N_PUBLISH_REMOTE)."; \
			exit 0; \
		fi; \
		echo "Module main is up-to-date; pushing missing tag $$tag for $$split_sha..."; \
	else \
		echo "==> Releasing i18n module tag $$tag to $(I18N_PUBLISH_REMOTE) (subtree $$split_sha)..."; \
		$(MAKE) i18n-publish; \
	fi; \
	if git rev-parse -q --verify "refs/tags/$$tag" >/dev/null; then \
		echo "Local tag $$tag already exists."; \
	else \
		git tag "$$tag" "$$split_sha"; \
	fi; \
	$(MAKE) -C $(I18N_TOOL_DIR) release-tag TAG="$$tag" REMOTE="$(I18N_PUBLISH_REMOTE)"

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
release-preflight:
	$(call CHECK_RELEASE_READY)

release-i18n-module:
	@if [ "$(RELEASE_I18N_MODULE)" = "1" ]; then \
		$(MAKE) i18n-module-release; \
	else \
		echo "==> Skipping i18n module release (RELEASE_I18N_MODULE=$(RELEASE_I18N_MODULE))"; \
	fi

release-push:
	@echo "==> Releasing ${APP_NAME} version $(VERSION)..."
	git push origin $(RELEASE_BRANCH)
	git push origin $(VERSION)

release:
	@$(MAKE) release-preflight
	@$(MAKE) release-i18n-module RELEASE_I18N_MODULE=$(RELEASE_I18N_MODULE)
	@$(MAKE) release-push

# Prepare for release (run on release branch before tagging)
prepare-release:
	$(call CHECK_ON_RELEASE_BRANCH,prepare-release)
	$(call CHECK_NO_VERSION_TAG)
	@$(MAKE) fmt-check
	@$(MAKE) i18n-sync
