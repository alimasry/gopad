build:
	@go build -o bin/gopad cmd/app/main.go

run: build
	@./bin/gopad

test:
	@go test -v ./test/...