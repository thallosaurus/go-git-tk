all: build

prereq:
	rm -rf build/*

build: prereq
	go build -o build/gittk-shell cmd/shell/main.go
	go build -o build/gittk-access cmd/access/main.go

install: build
	go install
