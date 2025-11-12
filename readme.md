# grpc golang topics 
- client side streaming
- server side streaming
- bidirectional streaming 

send metadata from client to server
md := metadata.Pairs("authorization", "Bearer=jwt-token")
ctx = metadata.NewOutgoingContext(ctx, md)
