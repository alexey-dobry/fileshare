package utils

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("x-user-id")
	if len(values) == 0 {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	return uuid.Parse(values[0])
}
