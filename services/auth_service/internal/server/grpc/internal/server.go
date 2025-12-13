package internal

import (
	pb "github.com/alexey-dobry/fileshare/pkg/gen/auth/intauth"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/jwt"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
)

type InternalServer struct {
	pb.UnimplementedAuthServer

	logger     logger.Logger
	store      store.Store
	jwtHandler jwt.JWTHandler
}

func New(logger logger.Logger, store store.Store, jwtHandler jwt.JWTHandler) *InternalServer {
	return &InternalServer{
		store:      store,
		logger:     logger.WithFields("layer", "grpc server api", "internal"),
		jwtHandler: jwtHandler,
	}
}
