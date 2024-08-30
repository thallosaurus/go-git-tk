CWD := $(shell pwd)

all: build

prereq:
	rm -rf build/*

build: prereq
	GOOS=darwin go build -o build/gittk.mac cmd/shell/main.go
	GOOS=windows go build -o build/gittk.windows.exe cmd/shell/main.go
	GOOS=linux go build -o build/gittk.linux cmd/shell/main.go
	GOOS=linux GOARCH=arm build -o build/gittk.rpi.linux cmd/shell/main.go

	GOOS=darwin go build -o build/gittk-access.mac cmd/access/main.go
	GOOS=windows go build -o build/gittk-access.windows.exe cmd/access/main.go
	GOOS=linux go build -o build/gittk-access.linux cmd/access/main.go
	GOOS=linux GOARCH=arm go build -o build/gittk-access.rpi.linux cmd/access/main.go

setup_tests: build
	git init --bare $(CWD)/repos/__testing_single_commit.git

test: setup_tests
	ROOT_CWD=$(CWD) go test ./...

clean_tests:
	rm -rf $(CWD)/repos/__testing_*.git

install: build
	go install
