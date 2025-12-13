package internal

import (
	"context"
	"strings"

	pb "github.com/alexey-dobry/fileshare/pkg/gen/auth/intauth"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/model"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/utils"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *InternalServer) RegisterCredentials(ctx context.Context, req *pb.RegisterCredentialsRequest) (*emptypb.Empty, error) {
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Errorf("Failed to generate password hash: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	newUserCredentials := model.Credentials{
		UUID:         uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
	}

	err = s.store.User().Add(newUserCredentials)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil, status.Error(codes.AlreadyExists, "Account with specified email already exists")
	} else if err != nil {
		s.logger.Errorf("Failed to add new user to data: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *InternalServer) DeleteCredentials(ctx context.Context, req *pb.DeleteCredentialsRequest) (*emptypb.Empty, error) {
	err := s.store.User().Delete(req.Email)
	if err != nil {
		s.logger.Errorf("Failed to delete user from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}
