PROJECT_NAME := gitlab-hook
BINARY_NAME := ${PROJECT_NAME}d
REPO_NAME := "gitlab.com/nerzhul/${PROJECT_NAME}"
PKG_LIST := $(shell go list ${REPO_NAME}/... | grep -v /vendor/)
BUILD_LD_FLAGS := $(shell echo "-X ${REPO_NAME}/cmd/${BINARY_NAME}/internal.AppBuildDate='`date -u '+%Y-%m-%d_%I:%M:%S%p'`' -X ${REPO_NAME}/cmd/${BINARY_NAME}/internal.AppVersion='`git describe --tags`'")

.PHONY: all dep doc build test

all: test lint doc build

lint: ## test
	@go get -u github.com/golang/lint/golint
	@${GOPATH}/bin/golint -set_exit_status ${PKG_LIST}

test: dep ## Run unittests
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

dep:
	@go get -u github.com/tools/godep
	@godep restore

gitlab-hook: dep
	@cd cmd/gitlab-hookd && \
    		mkdir -p ${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/ && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/gitlab-hookd"

slackbot: dep
	@cd cmd/slackbot && \
    		mkdir -p ${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/ && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/slackbot"

build: gitlab-hook slackbot

doc: swagger_doc

swagger_doc:
	@cd cmd/${BINARY_NAME} && \
		mkdir -p ${CI_PROJECT_DIR}/artifacts && \
		go get -u github.com/go-swagger/go-swagger/cmd/swagger && \
		${GOPATH}/bin/swagger generate spec -o ${CI_PROJECT_DIR}/artifacts/swagger.json