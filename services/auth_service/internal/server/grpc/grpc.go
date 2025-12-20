package grpc

import (
	intauthrpc "github.com/alexey-dobry/fileshare/pkg/gen/auth/intauth"
	pubauthrpc "github.com/alexey-dobry/fileshare/pkg/gen/auth/pubauth"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/jwt"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/server/grpc/internal"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/server/grpc/public"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewPublicServer(logger logger.Logger, repository store.Store, jwtHandler jwt.JWTHandler) *grpc.Server {
	s := grpc.NewServer()

	pubauthrpc.RegisterAuthServer(s, public.New(logger, repository, jwtHandler))

	reflection.Register(s)

	return s
}

func NewInternalServer(logger logger.Logger, repository store.Store, jwtHandler jwt.JWTHandler) *grpc.Server {
	s := grpc.NewServer()

	intauthrpc.RegisterInternalAuthServer(s, internal.New(logger, repository, jwtHandler))

	reflection.Register(s)

	return s
}
