package handlers

import mainapi "github.com/grpc_tutorials/my_project/proto/gen"

type Server struct {
	mainapi.UnimplementedExecsServiceServer
	mainapi.UnimplementedStudentsServiceServer
	mainapi.UnimplementedTeachersServiceServer
}
