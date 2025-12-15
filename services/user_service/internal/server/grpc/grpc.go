package grpc

import (
	intuserrpc "github.com/alexey-dobry/fileshare/pkg/gen/user/intuser"
	pubuserrpc "github.com/alexey-dobry/fileshare/pkg/gen/user/pubuser"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/server/grpc/internal"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/server/grpc/public"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store"

	"google.golang.org/grpc"
)

func NewPublicServer(logger logger.Logger, repository store.Store) *grpc.Server {
	s := grpc.NewServer()

	pubuserrpc.RegisterUserServer(s, public.New(logger, repository))

	return s
}

func NewInternalServer(logger logger.Logger, repository store.Store) *grpc.Server {
	s := grpc.NewServer()

	intuserrpc.RegisterUserServer(s, internal.New(logger, repository))

	return s
}
