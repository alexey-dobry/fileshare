package public

import (
	pb "github.com/alexey-dobry/fileshare/pkg/gen/file/pubfile"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store"
)

type PublicServer struct {
	pb.UnimplementedFileServiceServer

	logger logger.Logger
	store  store.Store
}

func New(logger logger.Logger, store store.Store) *PublicServer {
	return &PublicServer{
		store:  store,
		logger: logger.WithFields("layer", "grpc server api", "public"),
	}
}
