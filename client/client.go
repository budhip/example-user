package client

import (
	"context"
	"fmt"
	protoc "github.com/budhip/example-user/delivery/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	conf "github.com/budhip/example-user/config"
)

type userClient struct {
	conn   *grpc.ClientConn
	client protoc.UserHandlerClient
}

// UserClient to access user service
type UserClient interface {
	GetUserByID(ctx context.Context, userID int64) (*protoc.CreateUserResponse, error)
	Close() error
}

func (u userClient) Close() error {
	return u.conn.Close()
}

// NewUserClient :
func NewUserClient(address string) (UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := protoc.NewUserHandlerClient(conn)

	return &userClient{conn: conn, client: client}, nil
}

func (u userClient) GetUserByID(ctx context.Context, userID int64) (*protoc.CreateUserResponse, error) {
	fmt.Println("aaa")
	logger := conf.Logger(ctx)

	logger.Info("START GetUserByID in client")
	//conf.WithContext(context.Background()).Info("START GetUserByID in client")

	resp, err := u.client.GetUserByID(ctx, &protoc.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return &protoc.CreateUserResponse{
		Id:        resp.Id,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}
