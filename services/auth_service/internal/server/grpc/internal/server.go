package internal

import (
	pb "github.com/alexey-dobry/fileshare/pkg/gen/auth/intauth"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
)

type InternalServer struct {
	pb.UnimplementedInternalAuthServer

	logger logger.Logger
	store  store.Store
}

func New(logger logger.Logger, store store.Store) *InternalServer {
	return &InternalServer{
		store:  store,
		logger: logger.WithFields("layer", "grpc server api", "internal"),
	}
}
