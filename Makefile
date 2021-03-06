PROJECT_NAME := bot
REPO_NAME := "gitlab.com/nerzhul/${PROJECT_NAME}"
PKG_LIST := $(shell go list ${REPO_NAME}/... | grep -v /vendor/)
# BUILD_LD_FLAGS := $(shell echo "-X ${REPO_NAME}/cmd/${BINARY_NAME}/internal.AppBuildDate='`date -u '+%Y-%m-%d_%I:%M:%S%p'`' -X ${REPO_NAME}/cmd/${BINARY_NAME}/internal.AppVersion='`git describe --tags`'")

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
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
	mkdir -p ${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/

commandhandler: dep
	@cd cmd/commandhandler && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/commandhandler"

webhook: dep
	@cd cmd/webhookd && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/webhook"

ircbot: dep
	@cd cmd/ircbot && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/ircbot"

matterbot: dep
	@cd cmd/matterbot && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/matterbot"

releasechecker: dep
	@cd cmd/releasechecker && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/releasechecker"

slackbot: dep
	@cd cmd/slackbot && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/slackbot"

twitterbot: dep
	@cd cmd/twitterbot && \
    		go build  -ldflags "${BUILD_LD_FLAGS}" -o "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/twitterbot"

build: commandhandler webhook ircbot matterbot releasechecker slackbot twitterbot

install: build
	install -d /usr/local/etc/rc.d
	install -d /usr/local/bin
	install -m 0755 "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/commandhandler" /usr/local/bin/commandhandler
	install -m 0755 res/freebsd/commandhandler.sh /usr/local/etc/rc.d/commandhandler
	install -m 0755 "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/webhook" /usr/local/bin/webhook
	install -m 0755 res/freebsd/webhook.sh /usr/local/etc/rc.d/webhook
	install -m 0755 "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/ircbot" /usr/local/bin/ircbot
	install -m 0755 cmd/ircbot/res/freebsd_ircbot.sh /usr/local/etc/rc.d/ircbot
	install -m 0755 "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/matterbot" /usr/local/bin/matterbot
	install -m 0755 res/freebsd/matterbot.sh /usr/local/etc/rc.d/matterbot
	install -m 0755 "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/slackbot" /usr/local/bin/slackbot
	install -m 0755 res/freebsd/slackbot.sh /usr/local/etc/rc.d/slackbot
	install -m 0755 "${CI_PROJECT_DIR}/artifacts/${GOOS}_${GOARCH}/twitterbot" /usr/local/bin/twitterbot
	install -m 0755 res/freebsd/twitterbot.sh /usr/local/etc/rc.d/twitterbot

doc: swagger_doc

swagger_doc:
	@cd cmd/webhookd && \
		mkdir -p ${CI_PROJECT_DIR}/artifacts && \
		go get -u github.com/go-swagger/go-swagger/cmd/swagger && \
		${GOPATH}/bin/swagger generate spec -o ${CI_PROJECT_DIR}/artifacts/swagger.json