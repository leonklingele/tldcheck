BINARY_OUT ?= ./tldcheck
CMD_DIR := ./
TLDS_FILE := ./tlds-alpha-by-domain.txt

TEST_COVERAGE_OUT := ./.gocoverage

.PHONY: all
all:
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) build

.PHONY: lint
lint:
	@golangci-lint run -v \
		./...

.PHONY: test-fast
test-fast:
	go test -v \
		-shuffle on \
		-failfast \
		./...

.PHONY: test
test: test-fast
	go test -v \
		-shuffle on \
		-vet=all \
		-race \
		-cover -covermode=atomic -coverprofile="${TEST_COVERAGE_OUT}" \
		./...

.PHONY: test-cover-open
test-cover-open: test
	go tool cover \
		-html="${TEST_COVERAGE_OUT}"

BUILD_FLAGS ?=
.PHONY: build
build:
	go build -v \
		-o "${BINARY_OUT}" \
		${BUILD_FLAGS} \
		"${CMD_DIR}"

DEBUG ?= true
.PHONY: run
run:
	DEBUG="${DEBUG}" air -c ./.air.toml

.PHONY: clean
clean:
	go clean -r -cache -testcache -modcache -fuzzcache

.PHONY: update-tlds
update-tlds:
	curl --silent -o "$(TLDS_FILE)" https://data.iana.org/TLD/tlds-alpha-by-domain.txt
