> protoc --go_out=. --go-grpc_out=. proto/main.proto
# multiple proto file resolve the imports 
> protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto proto/greeter.proto
> protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto proto/greeter.proto proto/farewell.proto
> go mod tidy
> brew install protobuf
> protoc --version
# generate certificate
> openssl genrsa -out key.pem 2048
> openssl req -new -key key.pem -out cert.csr
> openssl x509 -req -days 365 -in cert.csr -signkey key.pem -out cert.pem
> openssl x509 -text -noout -in cert.pem

# using cert.conf file to generate certificate 
> openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config cert.conf

issue fixing
```
protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto proto/greeter.proto proto/farewell.proto     
protoc-gen-go: Go package "/proto/gen" has inconsistent names farewellpb (farewell.proto) and mainpipb (greeter.proto)
--go_out: protoc-gen-go: Plugin failed with status code 1.
```

The error indicates that two .proto files are trying to use the same Go package path (gen) but with different package names (farewellpb vs mainpipb). This creates a conflict.

To fix this, you have two options:

Option 1: Use separate package paths (Recommended)

Update the go_package option to use a unique path for each proto file:

```
syntax = "proto3";
package farewell;
option go_package = "/proto/gen/farewell;farewellpb";
```

# update the greet file option
```
option go_package = "/proto/gen/greeter;mainpipb";
```