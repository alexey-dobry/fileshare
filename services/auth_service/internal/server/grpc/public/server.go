package public

import (
	pb "github.com/alexey-dobry/fileshare/pkg/gen/auth/pubauth"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/jwt"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
)

type PublicServer struct {
	pb.UnimplementedAuthServer

	logger     logger.Logger
	store      store.Store
	jwtHandler jwt.JWTHandler
}

func New(logger logger.Logger, store store.Store, jwtHandler jwt.JWTHandler) *PublicServer {
	return &PublicServer{
		store:      store,
		logger:     logger.WithFields("layer", "grpc server api", "public"),
		jwtHandler: jwtHandler,
	}
}
