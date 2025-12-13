package public

import (
	"context"
	"errors"
	"time"

	pb "github.com/alexey-dobry/fileshare/pkg/gen/auth/pubauth"
	"github.com/alexey-dobry/fileshare/pkg/validator"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/jwt"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/utils"
	"gorm.io/gorm"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *PublicServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validator.V.Var(req.Email, "not null"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid login arguments")
	}

	if err := validator.V.Var(req.Password, "not null"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid login arguments")
	}

	user, err := s.store.User().GetOneByMail(req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		s.logger.Errorf("Failed to get user data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, status.Error(codes.PermissionDenied, "Wrong password")
	}

	refreshToken, accessToken, err := s.jwtHandler.GenerateJWTPair(jwt.Claims{
		UserID: user.UUID,
	})

	if err != nil {
		s.logger.Errorf("Failed to generate token pair: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *PublicServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	claims, err := s.jwtHandler.ValidateJWT(req.RefreshToken, jwt.RefreshToken)
	if errors.Is(err, jwt.ErrJWTTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, "JWT token expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		s.logger.Errorf("Failed validate refresh token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	isLogout, err := s.store.Blacklist().IsSessionLoggedOut(claims.ID)
	if err != nil {
		s.logger.Errorf("Failed to check session state: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	} else if isLogout {
		return nil, status.Error(codes.PermissionDenied, "Session is expired")
	}

	user, err := s.store.User().GetOneByID(claims.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		s.logger.Errorf("Failed to get user data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	refreshToken, accessToken, err := s.jwtHandler.GenerateJWTPair(jwt.Claims{
		UserID: user.UUID,
	})
	if err != nil {
		s.logger.Errorf("Failed to generate token pair: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *PublicServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	accessClaims, err := s.jwtHandler.ValidateJWT(req.AccessToken, jwt.AccessToken)
	if errors.Is(err, jwt.ErrJWTTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, "JWT token expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		s.logger.Errorf("Failed validate access token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	refreshClaims, err := s.jwtHandler.ValidateJWT(req.RefreshToken, jwt.RefreshToken)
	if errors.Is(err, jwt.ErrJWTTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, "JWT token expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		s.logger.Errorf("Failed validate refresh token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	err = s.store.Blacklist().BlacklistAccessToken(accessClaims.ID, accessClaims.ExpiresAt.Sub(time.Now()))
	if err != nil {
		s.logger.Errorf("Failed to blacklist access token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	err = s.store.Blacklist().StoreLogoutSession(refreshClaims.ID, refreshClaims.ExpiresAt.Sub(time.Now()))
	if err != nil {
		s.logger.Errorf("Failed to blacklist refresh token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

func (s *PublicServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	claims, err := s.jwtHandler.ValidateJWT(req.AccessToken, jwt.AccessToken)
	if errors.Is(err, jwt.ErrJWTTokenExpired) {
		return nil, status.Error(codes.Unauthenticated, "JWT token expired")
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		s.logger.Errorf("Failed validate access token: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if err := validator.V.Var(req.OldPassword, "not null"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid login arguments")
	}

	user, err := s.store.User().GetOneByID(claims.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		s.logger.Errorf("Failed to get user data from database: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return nil, status.Error(codes.PermissionDenied, "Wrong password")
	}

	isLogout, err := s.store.Blacklist().IsAccessTokenBlacklisted(claims.ID)
	if err != nil {
		s.logger.Errorf("Failed to check session state: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	} else if isLogout {
		return nil, status.Error(codes.PermissionDenied, "Access token is blacklisted")
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Errorf("Failed to generate password hash: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	err = s.store.User().UpdatePassword(claims.UserID, newHash)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "User entry with given credentials not found")
	} else if err != nil {
		s.logger.Errorf("Failed to update user password hash: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}
