.DEFAULT_GOAL = build
SHELL = /bin/bash

RESET := $(shell tput sgr0)
BOLD := $(shell tput bold)
RED := $(shell tput setaf 1)
EOL := \n

setup-environment:
	@printf '${BOLD}${RED}make: *** [setup-environment]${RESET}${EOL}'
	$(eval SRC_PATH := $(shell pwd))
	$(eval TMP_PATH := $(shell mktemp -d))
	$(eval GIT_PATH := $(shell go list))
	$(eval MOD_PATH := ${TMP_PATH}/src/${GIT_PATH})
	$(eval VER_DESC := $(shell git describe --tags))
	export GOPATH='${TMP_PATH}'
	export GO111MODULE='on'
	mkdir -p ${MOD_PATH}
	rmdir ${MOD_PATH}
	ln -Fs ${SRC_PATH} ${MOD_PATH}
.PHONY: setup-environment

build: plugin
.PHONY: build

plugin: setup-environment
	@printf '${BOLD}${RED}make: *** [plugin]${RESET}${EOL}'
	cd ${MOD_PATH}                              && \
	go build                                       \
		-o 'plugin'                                \
		-a                                         \
		-installsuffix 'cgo'                       \
		-gcflags 'all=-trimpath "${TMP_PATH}/src"' \
		-v

test: setup-environment
	@printf '${BOLD}${RED}make: *** [test]${RESET}${EOL}'
	cd ${MOD_PATH} && \
	go test           \
		-v ./...
.PHONY: test

continuous-integration: test build
.PHONY: continuous-integration
