TLDS_FILE := ./cmd/tldcheck/tlds-alpha-by-domain.txt

.PHONY: all
all: build

.PHONY: build
build:
	go build ./...

.PHONY: install
install: build
	go install ./...

.PHONY: update-tlds
update-tlds:
	curl --silent -o $(TLDS_FILE) https://data.iana.org/TLD/tlds-alpha-by-domain.txt
