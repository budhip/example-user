package grpc

import (
	"context"
	"github.com/budhip/example-user/delivery/pb"

	"github.com/budhip/example-user/model"
)

func (s *server) CreateUser(ctx context.Context,
	request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	reqIns := &model.AddNewUserRequest{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
	}
	resp, err := s.service.CreateUser(ctx, reqIns)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}

func (s *server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.CreateUserResponse, error) {
	resp, err := s.service.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id:        resp.ID,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}
