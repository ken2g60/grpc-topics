protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto proto/students.proto proto/execs.proto

# start & stop mongodb service
> brew services start mongodb-community
> brew services stop mongodb-community

