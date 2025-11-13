protoc --go_out=. main.proto
protoc --go_out=. --go-grpc=. proto/main.proto
protoc -I=. --go_out=.. --go-grpc_out=.. --validate_out=lang=go:.. main.proto

protoc-gen-validate
builds validation whiles building the proto.