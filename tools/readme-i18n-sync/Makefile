SOURCE ?= README.md
I18N_DIR ?= i18n
TM_DIR ?=
I18N_FORCE ?=
REMOTE ?= origin
TAG_PREFIX ?= readme-i18n-sync/v
GOFLAGS ?= -mod=vendor
NEXT_TAG = $(shell last=$$(git tag --list '$(TAG_PREFIX)*' | sed 's#$(TAG_PREFIX)##' | sort -V | tail -n1); if [ -z "$$last" ]; then echo '$(TAG_PREFIX)0.1.0'; else echo "$$last" | awk -F. '{printf "$(TAG_PREFIX)%d.%d.%d", $$1, $$2, $$3+1}'; fi)
TAG ?= $(NEXT_TAG)

GOFMT_PATHS = -path ./vendor -o -path ./.git -o -path ./.cache -o -path ./bin
RUN_BASE = go run ./cmd/readme-i18n-sync --source $(SOURCE) --i18n-dir $(I18N_DIR) $(if $(TM_DIR),--tm-dir $(TM_DIR),)

.PHONY: test fmt fmt-check check update sync vendor vendor-check print-next-tag tag release-tag release

test:
	@echo "==> Testing readme-i18n-sync..."
	GOFLAGS="$(GOFLAGS)" go test ./...

fmt:
	@echo "==> Formatting readme-i18n-sync..."
	@find . \( $(GOFMT_PATHS) \) -prune -o -type f -name '*.go' -print0 | xargs -0r gofmt -w

fmt-check:
	@echo "==> Checking formatting for readme-i18n-sync..."
	@fmt_issues="$$(find . \( $(GOFMT_PATHS) \) -prune -o -type f -name '*.go' -print0 | xargs -0r gofmt -l)"; \
	if [ -n "$$fmt_issues" ]; then \
		echo "gofmt issues:"; \
		echo "$$fmt_issues"; \
		exit 1; \
	fi

check:
	@echo "==> Checking translations..."
	GOFLAGS="$(GOFLAGS)" $(RUN_BASE) --check

update:
	@echo "==> Updating translations..."
	GOFLAGS="$(GOFLAGS)" $(RUN_BASE) $(if $(I18N_FORCE),--force)

sync:
	@echo "==> Syncing translations..."
	@$(MAKE) update I18N_FORCE=$(I18N_FORCE) SOURCE="$(SOURCE)" I18N_DIR="$(I18N_DIR)" TM_DIR="$(TM_DIR)"
	@$(MAKE) check SOURCE="$(SOURCE)" I18N_DIR="$(I18N_DIR)" TM_DIR="$(TM_DIR)"

vendor:
	@echo "==> Syncing module dependencies and vendor..."
	GOFLAGS= go mod tidy
	GOFLAGS= go mod vendor

vendor-check:
	@echo "==> Verifying module vendor is up to date..."
	GOFLAGS= go mod tidy
	GOFLAGS= go mod vendor
	@if ! git diff --quiet -- go.mod go.sum vendor/modules.txt; then \
		echo "Vendor is out of date. Run: make vendor"; \
		git --no-pager diff -- go.mod go.sum vendor/modules.txt; \
		exit 1; \
	fi

print-next-tag:
	@echo "$(NEXT_TAG)"

tag:
	@echo "==> Creating module tag $(TAG)..."
	git tag $(TAG)

release-tag:
	@echo "==> Pushing module tag $(TAG) to $(REMOTE)..."
	git push $(REMOTE) $(TAG):refs/tags/$(TAG)

release:
	@$(MAKE) tag TAG="$(TAG)"
	@$(MAKE) release-tag TAG="$(TAG)" REMOTE="$(REMOTE)"
