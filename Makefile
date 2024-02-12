build:
	@go build -o ./bin/server ./cmd/server/main.go

run: build
	@./bin/server -http=:4000 -grpc=:4001 -keypath "jwt"

test:
	go test -v -cover ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/service.proto

.PHONY: proto
