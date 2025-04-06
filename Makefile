default: build

download:
	@go run cmd/download/main.go

lint:
	@golangci-lint run

clean:
	@rm -rf bin