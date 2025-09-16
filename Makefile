install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

build-swagger:
	swag init -g cmd/plata_currency_quotation/main.go

codegen:
	$(MAKE) build-swagger

lint:
	golangci-lint run

run-dev:
	go run cmd/plata_currency_quotation/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -o server cmd/plata_currency_quotation/main.go

initial-setup:
	go mod download
	$(MAKE) install-tools

test:
	go test -race -v ./...
