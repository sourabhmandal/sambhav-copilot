package user

import (
	"database/sql"
	"errors"
	"fmt"
	"sambhavhr/internal/repository"
	"strconv"

	"context"
)

type reportServiceSqlc struct {
	userRepository repository.Querier
}

func NewReportService(userRepo repository.Querier) ReportService {
	return &reportServiceSqlc{userRepository: userRepo}
}

// RegisterUser registers a new user in the system.
func (u *reportServiceSqlc) RegisterUser(ctx context.Context, name, email string) error {
	// Check if the user already exists based on email
	existingUser, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		// Create a new user in the schema
		_, err = u.userRepository.CreateUser(ctx, &repository.CreateUserParams{
			Email: email,
			Name:  name,
			Bio:   nil, // Assuming empty bio for new user
		})
		if err != nil {
			return err
		}
	} else if existingUser != nil {
		return errors.New("user already exists")
	} else {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by their ID.
func (u *reportServiceSqlc) GetAllUsers(ctx context.Context) ([]*repository.User, error) {

	// Get the user by ID from the schema
	users, err := u.userRepository.ListAllUsers(ctx)
	fmt.Println(users, err)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID retrieves a user by their ID.
func (u *reportServiceSqlc) GetUserByID(ctx context.Context, userID string) (*repository.User, error) {
	// Convert userID to int64 if necessary
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get the user by ID from the schema
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
