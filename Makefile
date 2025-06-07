.PHONY: proto clean install-tools all

PROTO_DIR=api/study

proto:
	protoc \
		--proto_path=api/study \
		--go_out=api/study --go_opt=paths=source_relative \
		--go-grpc_out=api/study --go-grpc_opt=paths=source_relative \
		api/study/*.proto

clean:
	rm -f $(PROTO_DIR)/*.pb.go
	rm -f $(PROTO_DIR)/*_grpc.pb.go

install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

all: proto
