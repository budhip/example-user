package service

import (
	"context"
	"time"

	"github.com/budhip/example-user/model"
	"github.com/budhip/example-user/repository"

	conf "github.com/budhip/example-user/config"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.AddNewUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

type userService struct {
	userRepo      repository.PostgreUserRepository
	snowFlake     conf.Snowflake
}

// NewUserService will create new an userService object representation of UserService interface
func NewUserService(m repository.PostgreUserRepository, sf conf.Snowflake) UserService {
	return &userService{m, sf}
}

func (us *userService) CreateUser(ctx context.Context, u *model.AddNewUserRequest) (*model.User, error) {
	idSnow := us.snowFlake.GetNewID()

	userID := idSnow
	createdAt := time.Now().Format("2006-01-02 15:04:5")
	updatedAt := time.Now().Format("2006-01-02 15:04:5")

	userModel := &model.User{
		ID:        userID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := us.userRepo.StoreUser(userModel)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

func (us *userService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	logger := conf.Logger(ctx)

	logger.Info("START GetUserByID")

	list, err := us.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return list, nil
}