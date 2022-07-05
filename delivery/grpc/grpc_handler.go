package grpc

import (
	"github.com/budhip/example-user/delivery/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userSrv "github.com/budhip/example-user/service"
)

type server struct {
	service userSrv.UserService
}

func NewUserServerGRPC(gServer *grpc.Server, userSrv userSrv.UserService) {
	// RegisterUserHandlerServer using Server struct
	userServer := &server{
		service: userSrv,
	}
	pb.RegisterUserHandlerServer(gServer, userServer)

	reflection.Register(gServer)
}
