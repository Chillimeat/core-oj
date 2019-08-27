

JUDGER_CONTAINTER_NAME=core-oj/judger

judger: build
	docker build --tag ${JUDGER_CONTAINTER_NAME} judgerx

build:
	env GOOS=linux GOARCH=386 go build -o ./judgerx/judger  ./judgerx/main/
	g++ ./std-checker/stdin-outfile-stderr.cpp -o ./judgerx/sos-checker
	chmod +x ./judgerx/judger
	chmod +x ./judgerx/sos-checker

test:
	./judgerx/judger ./test


.PHONY: judger build test





