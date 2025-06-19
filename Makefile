APP_NAME ?= hikvision-backup
RELEASE_BRANCH ?= master

# Получаем текущую версию из тега git
VERSION ?= $(shell git describe --tags --abbrev=0)

# Удаляем префикс 'v' если он есть
RAW_VERSION := $(subst v,,$(VERSION))

# Разделяем версию на части (major.minor.patch)
MAJOR := $(word 1, $(subst ., ,$(RAW_VERSION)))
MINOR := $(word 2, $(subst ., ,$(RAW_VERSION)))
PATCH := $(word 3, $(subst ., ,$(RAW_VERSION)))

# Увеличиваем patch-версию на единицу
NEW_PATCH := $(shell echo $$(( $(PATCH) + 1 )))
NEW_VERSION_PATCH := v$(MAJOR).$(MINOR).$(NEW_PATCH)

# Увеличиваем minor-версию и сбрасываем patch-версию до 0
NEW_MINOR := $(shell echo $$(( $(MINOR) + 1 )))
NEW_VERSION_MINOR := v$(MAJOR).$(NEW_MINOR).0

# Увеличиваем major-версию и сбрасываем minor и patch версии до 0
NEW_MAJOR := $(shell echo $$(( $(MAJOR) + 1 )))
NEW_VERSION_MAJOR := v$(NEW_MAJOR).0.0

NEW_VERSION_ALPHA := $(NEW_VERSION_PATCH)-alpha
NEW_VERSION_BETA := $(NEW_VERSION_PATCH)-beta

COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS_STRING = -X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.buildDate=${DATE}'
LDFLAGS = -ldflags="${LDFLAGS_STRING}"

.PHONY: build
build:
	@echo "==> Building ${APP_NAME}..."
	go build ${LDFLAGS} -o ${APP_NAME} .

next-alpha:
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_ALPHA)..."
	git tag $(NEW_VERSION_ALPHA)

next-beta:
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_BETA)..."
	git tag $(NEW_VERSION_BETA)

next-patch:
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_PATCH)..."
	git tag $(NEW_VERSION_PATCH)

next-minor:
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_MINOR)..."
	git tag $(NEW_VERSION_MINOR)

next-major:
	@echo "==> New ${APP_NAME} version $(NEW_VERSION_MAJOR)..."
	git tag $(NEW_VERSION_MAJOR)

release:
	@echo "==> Releasing ${APP_NAME} version $(VERSION)..."
	git push origin $(RELEASE_BRANCH)
	git push origin $(VERSION)
