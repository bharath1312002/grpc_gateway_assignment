package service

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"log"
	"user_service/protogen/user"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	session *gocql.Session
}

func NewUserServiceServer(session *gocql.Session) *UserServiceServer {
	return &UserServiceServer{session: session}
}

// CreateUser creates a new user in the database.
func (s *UserServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	id := uuid.New().String() // Generate a unique ID for the user
	query := `INSERT INTO users (id, first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if err := s.session.Query(query, id, req.FirstName, req.LastName, req.Gender, req.DateOfBirth, req.PhoneNumber, req.Email, false).Exec(); err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}
	return &user.UserResponse{
		User: &user.User{
			Id:          id,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      req.Gender,
			DateOfBirth: req.DateOfBirth,
			PhoneNumber: req.PhoneNumber,
			Email:       req.Email,
			IsBlocked:   false,
		},
	}, nil
}

// UpdateUser updates an existing user's details.
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	query := `UPDATE users SET first_name = ?, last_name = ?, gender = ?, date_of_birth = ? WHERE id = ?`
	if err := s.session.Query(query, req.FirstName, req.LastName, req.Gender, req.DateOfBirth, req.Id).Exec(); err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}

	// Fetch the updated user details
	var (
		firstName, lastName, gender, dateOfBirth, phoneNumber, email string
		isBlocked                                                    bool
	)
	if err := s.session.Query(`SELECT first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked FROM users WHERE id = ?`, req.Id).Scan(
		&firstName, &lastName, &gender, &dateOfBirth, &phoneNumber, &email, &isBlocked,
	); err != nil {
		log.Printf("Failed to fetch updated user: %v", err)
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:          req.Id,
			FirstName:   firstName,
			LastName:    lastName,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			PhoneNumber: phoneNumber,
			Email:       email,
			IsBlocked:   isBlocked,
		},
	}, nil
}

// BlockUser blocks a user by setting the is_blocked flag to true.
func (s *UserServiceServer) BlockUser(ctx context.Context, req *user.BlockUserRequest) (*user.UserResponse, error) {
	query := `UPDATE users SET is_blocked = true WHERE id = ?`
	if err := s.session.Query(query, req.Id).Exec(); err != nil {
		log.Printf("Failed to block user: %v", err)
		return nil, err
	}

	// Fetch the updated user details
	var (
		firstName, lastName, gender, dateOfBirth, phoneNumber, email string
		isBlocked                                                    bool
	)
	if err := s.session.Query(`SELECT first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked FROM users WHERE id = ?`, req.Id).Scan(
		&firstName, &lastName, &gender, &dateOfBirth, &phoneNumber, &email, &isBlocked,
	); err != nil {
		log.Printf("Failed to fetch blocked user: %v", err)
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:          req.Id,
			FirstName:   firstName,
			LastName:    lastName,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			PhoneNumber: phoneNumber,
			Email:       email,
			IsBlocked:   isBlocked,
		},
	}, nil
}

// UnblockUser unblocks a user by setting the is_blocked flag to false.
func (s *UserServiceServer) UnblockUser(ctx context.Context, req *user.UnblockUserRequest) (*user.UserResponse, error) {
	query := `UPDATE users SET is_blocked = false WHERE id = ?`
	if err := s.session.Query(query, req.Id).Exec(); err != nil {
		log.Printf("Failed to unblock user: %v", err)
		return nil, err
	}

	// Fetch the updated user details
	var (
		firstName, lastName, gender, dateOfBirth, phoneNumber, email string
		isBlocked                                                    bool
	)
	if err := s.session.Query(`SELECT first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked FROM users WHERE id = ?`, req.Id).Scan(
		&firstName, &lastName, &gender, &dateOfBirth, &phoneNumber, &email, &isBlocked,
	); err != nil {
		log.Printf("Failed to fetch unblocked user: %v", err)
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:          req.Id,
			FirstName:   firstName,
			LastName:    lastName,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			PhoneNumber: phoneNumber,
			Email:       email,
			IsBlocked:   isBlocked,
		},
	}, nil
}

// UpdateContact updates a user's phone number and/or email.
func (s *UserServiceServer) UpdateContact(ctx context.Context, req *user.UpdateUserContactRequest) (*user.UserResponse, error) {
	query := `UPDATE users SET phone_number = ?, email = ? WHERE id = ?`
	if err := s.session.Query(query, req.PhoneNumber, req.Email, req.Id).Exec(); err != nil {
		log.Printf("Failed to update contact: %v", err)
		return nil, err
	}

	// Fetch the updated user details
	var (
		firstName, lastName, gender, dateOfBirth, phoneNumber, email string
		isBlocked                                                    bool
	)
	if err := s.session.Query(`SELECT first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked FROM users WHERE id = ?`, req.Id).Scan(
		&firstName, &lastName, &gender, &dateOfBirth, &phoneNumber, &email, &isBlocked,
	); err != nil {
		log.Printf("Failed to fetch updated user: %v", err)
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:          req.Id,
			FirstName:   firstName,
			LastName:    lastName,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			PhoneNumber: phoneNumber,
			Email:       email,
			IsBlocked:   isBlocked,
		},
	}, nil
}

// GetUser retrieves a user by phone number or email.
func (s *UserServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	var (
		id, firstName, lastName, gender, dateOfBirth, phoneNumber, email string
		isBlocked                                                        bool
	)

	query := `SELECT id, first_name, last_name, gender, date_of_birth, phone_number, email, is_blocked FROM users WHERE phone_number = ? OR email = ? LIMIT 1`
	if err := s.session.Query(query, req.PhoneNumber, req.Email).Scan(
		&id, &firstName, &lastName, &gender, &dateOfBirth, &phoneNumber, &email, &isBlocked,
	); err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:          id,
			FirstName:   firstName,
			LastName:    lastName,
			Gender:      gender,
			DateOfBirth: dateOfBirth,
			PhoneNumber: phoneNumber,
			Email:       email,
			IsBlocked:   isBlocked,
		},
	}, nil
}
