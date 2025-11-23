package handlers

import (
	"context"
	"fmt"

	mainapi "github.com/grpc_tutorials/my_project/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	mainapi.UnimplementedExecsServiceServer
	mainapi.UnimplementedStudentsServiceServer
	mainapi.UnimplementedTeachersServiceServer
}

func (s Server) GetExecs(ctx context.Context, req *mainapi.GetExecRequest) (*mainapi.Execs, error) {
	fmt.Println(req.Exec)
	return nil, status.Errorf(codes.Unimplemented, "method GetExecs not implemented")
}
