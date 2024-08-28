all: build

prereq:
	rm -rf build/*

build: prereq
	go build -o build/shell cmd/shell/main.go
	go build -o build/access cmd/access/main.go

install: build
	go install
