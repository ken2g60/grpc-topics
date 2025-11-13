protoc --go_out=. main.proto
protoc --go_out=. --go-grpc=. proto/main.proto
protoc -I=. --go_out=.. --go-grpc_out=.. --grpc-gateway_out=.. --validate_out=lang=go:.. main.proto

protoc-gen-validate
builds validation whiles building the proto.

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

# fixing the annotation issues proto works
protoc -I=. --go_out=.. --go-grpc_out=.. --grpc-gateway_out=.. --validate_out=lang=go:. main.proto


wget -O annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/annotations.proto
wget -O http.proto https://raw.githubusercontent.com/googleapis/googleapis/refs/heads/master/google/api/http.proto


protoc -I. \
  -I/usr/local/include \
  -I$(go env GOPATH)/src \
  -I$(go env GOPATH)/github.com/grpc_tutorials/protoc-validate-restapi/google/api/ \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
  main.proto