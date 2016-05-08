all: build install 

build:
	go build

install:
	go install

test: build install
	go test -v *.go
