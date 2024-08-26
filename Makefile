prereq:
	rm -rf build/git-tk-shell

build: prereq
	go build git-tk-shell.go

install: build
	go install

all: build install