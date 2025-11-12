protoc \
--proto_path=protobuf "protobuf/service.proto" \
		--go_out=services/common/genproto/tutorials --go_opt=paths=source_relative \
  	--go-grpc_out=services/common/genproto/tutorials --go-grpc_opt=paths=source_relative

> protoc --go_out=. --go-grpc_out=. proto/main.proto # grpc-topics
