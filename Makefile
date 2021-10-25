.PHONY: mocks
mocks:
	@rm -rf ./internal/mocks/
	mockgen -destination=./internal/mocks/repo_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/repo EventRepo
	mockgen -destination=./internal/mocks/sender_mock.go -package=mocks github.com/ozonmp/est-water-api/internal/app/sender EventSender

.PHONY: run
run:
	go run cmd/est-water-api/main.go

.PHONY: test
test:
	make mocks
	go clean -testcache
	go test ./internal/app/retranslator -v
	go test ./internal/app/consumer -v
	go test ./internal/app/producer -v

.PHONY: build
build:
	go build -o bot cmd/est-water-api/main.go