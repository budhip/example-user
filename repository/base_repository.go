package repository

import (
	"database/sql"

	"github.com/budhip/example-user/model"
)

type postgreRepository struct {
	Conn *sql.DB
}

func NewPostgreRepository(conn *sql.DB) PostgreUserRepository {
	return &postgreRepository{conn}
}

type PostgreUserRepository interface {
	StoreUser(a *model.User) error
	GetUserByID(id int64) (*model.User, error)
}

type RedisRepository interface {
}
