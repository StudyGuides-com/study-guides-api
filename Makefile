.PHONY: proto clean install-tools build run all run-dev run-test run-prod

PROTO_DIR=api/study
BIN_DIR=bin
CMD_DIR=cmd/server

proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/search/tag_path.proto \
		$(PROTO_DIR)/search/search.proto \
		$(PROTO_DIR)/health/health.proto \
		$(PROTO_DIR)/study.proto

build:
	go build -o $(BIN_DIR)/study-server ./$(CMD_DIR)

run: build
	./$(BIN_DIR)/study-server

run-dev: build
	cp .env.dev .env
	./$(BIN_DIR)/study-server

run-test: build
	cp .env.test .env
	./$(BIN_DIR)/study-server

run-prod: build
	cp .env.prod .env
	./$(BIN_DIR)/study-server

clean:
	rm -f $(PROTO_DIR)/*.pb.go
	rm -f $(PROTO_DIR)/*/*.pb.go
	rm -f $(PROTO_DIR)/*_grpc.pb.go
	rm -rf $(BIN_DIR)
	rm -f .env

fmt:
	find . -type f -name '*.go' ! -name '*.pb.go' -exec gofmt -s -w {} +

lint:
	golangci-lint run --timeout=2m --skip-dirs-use-default --skip-files='.*\.pb\.go'


evans-dev:
	cp .env.dev .env
	evans --host localhost --port 1973 --header "Authorization=Bearer $$(cat .jwt.dev)" -r

evans-test:
	cp .env.test .env
	evans --host localhost --port 1973 --header "Authorization=Bearer $$(cat .jwt.test)" -r

evans-prod:
	cp .env.prod .env
	evans --host localhost --port 1973 --header "Authorization=Bearer $$(cat .jwt.prod)" -r

generate-tokens:
	node scripts/generate_jwt.js > .jwt.dev
	node scripts/generate_jwt.js > .jwt.test
	node scripts/generate_jwt.js > .jwt.prod

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

all: proto build
