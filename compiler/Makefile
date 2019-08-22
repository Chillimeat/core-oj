

COMPILER_CONTAINTER_NAME=core-oj/compiler

compiler: build
	docker build --tag ${COMPILER_CONTAINTER_NAME} compilerx

build:
	env GOOS=linux GOARCH=386 go build -o ./compilerx/compiler  ./compilerx/main/
	chmod +x ./compilerx/compiler

test:
	go test -v ./compilerx/


.PHONY: compiler build test





