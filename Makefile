default: build

download:
	@go run cmd/download/main.go

count:
	@go run cmd/count/main.go

lint:
	@golangci-lint run

clean:
	@rm -rf bin