.PHONY: install test cover run build start clean

install:
    go mod download

test:
	@echo "Run unit testing ..."
		mkdir -p ./coverage && \
			go test -v -coverprofile=./coverage/coverage.out -covermode=atomic ./bin/modules/user/usecases/...

cover: test
	@echo "Generating coverprofile ..."
		go tool cover -func=./coverage/coverage.out &&\
			go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

run:
    go run ./bin/app/main.go

build:
    go build -tags musl -o main ./bin/app

start:
    ./main

clean:
    @echo "Cleansing the last built ..."
    rm -rf bin