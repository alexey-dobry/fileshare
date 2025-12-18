package grpc

import (
	pb "github.com/alexey-dobry/fileshare/pkg/gen/file/pubfile"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/server/grpc/public"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store"

	"google.golang.org/grpc"
)

func NewPublicServer(logger logger.Logger, repository store.Store) *grpc.Server {
	s := grpc.NewServer()

	pb.RegisterFileServiceServer(s, public.New(logger, repository))
	return s
}
