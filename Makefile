all: build

prereq:
	rm -rf build/git-tk-shell

build: prereq
	go build -o build/git-tk-shell git-tk-shell.go

install: build
	go install
