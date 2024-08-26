prereq:
	rm -rf build/main

build: prereq
	go build -o build/main cmd/main/main.go

all: build