include .env.local
export
gen-proto:
	protoc -I ./pkg/api/proto/ \
	--go-grpc_out=pkg/api/go/ \
	--go_out=pkg/api/go/ \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out pkg/api/go/ \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out pkg/api/openapi/ \
	pkg/api/proto/*.proto
gen-mocks:
	mockgen -source internal/service/service.go -destination internal/service/service_mock_test.go -package service
test:
	go test -v -race ./... -coverprofile=coverage.out
cover:
	go tool cover -html=coverage.out
build:
	go build ./cmd/...
database:
	docker compose up -d
run: build database
	./notpastebin