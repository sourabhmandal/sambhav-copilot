package user

import (
	"sambhavhr/internal/repository"

	"context"
)

type ReportService interface {
	RegisterUser(ctx context.Context, name, email string) error
	GetAllUsers(ctx context.Context) ([]*repository.User, error)
	GetUserByID(ctx context.Context, userID string) (*repository.User, error)
	// Other methods...
}
