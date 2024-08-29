CWD := $(shell pwd)

all: build

prereq:
	rm -rf build/*

build: prereq
	go build -o build/gittk-shell cmd/shell/main.go
	go build -o build/gittk-access cmd/access/main.go

setup_tests: build
	git init --bare $(CWD)/repos/__testing_single_commit.git

test: setup_tests
	ROOT_CWD=$(CWD) go test ./...

clean_tests:
	rm -rf $(CWD)/repos/__testing_*.git

install: build
	go install
