package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/pkg/logger/zap"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/config"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/jwt"
	authrpc "github.com/alexey-dobry/fileshare/services/auth_service/internal/server/grpc"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store/authdb"
	"google.golang.org/grpc"
)

type App interface {
	Run(context.Context)
}

type app struct {
	publicServer          *grpc.Server
	internalServer        *grpc.Server
	publicServerAddress   string
	internalServerAddress string
	store                 store.Store
	logger                logger.Logger
}

func New(cfg config.Config) App {
	var a app
	var err error

	a.logger = zap.NewLogger(cfg.Logger).WithFields("layer", "app")

	a.publicServerAddress = fmt.Sprintf(":%s", cfg.GRPC.PublicPort)
	a.internalServerAddress = fmt.Sprintf(":%s", cfg.GRPC.InternalPort)

	a.store, err = authdb.New(a.logger, cfg.Store)
	if err != nil {
		a.logger.Fatalf("Failed to create store instance: %s", err)
	}

	jwtHandler, err := jwt.NewHandler(cfg.JWT)
	if err != nil {
		a.logger.Fatalf("Failed to create jwt handler: %s", err)
	}

	a.internalServer = authrpc.NewInternalServer(a.logger, a.store, jwtHandler)
	a.publicServer = authrpc.NewPublicServer(a.logger, a.store, jwtHandler)

	a.logger.Info("app was built")
	return &a
}

func (a *app) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	pubListener, err := net.Listen("tcp", a.publicServerAddress)
	if err != nil {
		a.logger.Fatal(err)
	}

	intListener, err := net.Listen("tcp", a.internalServerAddress)
	if err != nil {
		a.logger.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		a.logger.Error("Starting public grpc server at address %s...", a.publicServerAddress)
		if err := a.publicServer.Serve(pubListener); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				a.logger.Errorf("Grpc server error: %s", err)
				cancel()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		a.logger.Error("Starting internal grpc server at address %s...", a.internalServerAddress)
		if err := a.internalServer.Serve(intListener); err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				a.logger.Errorf("Grpc server error: %s", err)
				cancel()
			}
		}
	}()
	a.logger.Info("App is running...")

	select {
	case <-quit:
		a.logger.Info("shutdown signal received")
	case <-ctx.Done():
		a.logger.Info("context canceled")
	}

	a.logger.Info("stopping all services")

	cancel()
	if err := intListener.Close(); err != nil {
		a.logger.Warnf("Internal net listener closing ended with error: %s", err)
	}

	if err := pubListener.Close(); err != nil {
		a.logger.Warnf("Public net listener closing ended with error: %s", err)
	}
	wg.Wait()

	if err := a.store.Close(); err != nil {
		a.logger.Warnf("store closing ended with error: %s", err)
	}

	a.logger.Info("app was gracefully shutdown")
}
