package repository

import (
	"context"

	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error

	FindByID(ctx context.Context, id string) (*entity.User, error)

	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	FindAll(ctx context.Context, skip, limit int64) ([]*entity.User, error)

	Count(ctx context.Context) (int64, error)

	Update(ctx context.Context, id string, updates bson.M) error

	Delete(ctx context.Context, id string) error
}
