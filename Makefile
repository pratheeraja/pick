.PHONY: install test build

all: build

install:
	gb vendor restore
	gb build cmd/pick
	#cp ./bin/pick /usr/local/bin

test:
	gb test

build:
	gb build cmd/pick
