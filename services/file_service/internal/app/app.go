package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/alexey-dobry/fileshare/pkg/gen/file/pubfile"
	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/pkg/logger/zap"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/config"
	rpc "github.com/alexey-dobry/fileshare/services/file_service/internal/server/grpc"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App interface {
	Run(context.Context)
}

type app struct {
	publicServer        *grpc.Server
	publicServerAddress string

	gatewayServer  *http.Server
	gatewayAddress string

	store  store.Store
	logger logger.Logger
}

func New(cfg config.Config) App {
	var a app
	var err error

	a.logger = zap.NewLogger(cfg.Logger).WithFields("layer", "app")

	a.publicServerAddress = fmt.Sprintf(":%s", cfg.GRPC.PublicPort)
	a.gatewayAddress = fmt.Sprintf(":%s", cfg.GRPC.GatewayPort)

	a.store, err = file.New(a.logger, cfg.Store)
	if err != nil {
		a.logger.Fatalf("Failed to create store instance: %s", err)
	}

	a.publicServer = rpc.NewPublicServer(a.logger, a.store)

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

	if err := a.initGateway(ctx); err != nil {
		a.logger.Fatalf("failed to init gateway: %s", err)
	}

	var wg sync.WaitGroup

	// gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Infof("Starting public grpc server at %s", a.publicServerAddress)
		if err := a.publicServer.Serve(pubListener); err != nil {
			if ctx.Err() == nil {
				a.logger.Errorf("grpc server error: %s", err)
				cancel()
			}
		}
	}()

	// HTTP gateway
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Infof("Starting http gateway at %s", a.gatewayAddress)
		if err := a.gatewayServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Errorf("gateway error: %s", err)
			cancel()
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

	_ = a.gatewayServer.Shutdown(context.Background())
	_ = pubListener.Close()

	wg.Wait()

	if err := a.store.Close(); err != nil {
		a.logger.Warnf("store closing ended with error: %s", err)
	}

	a.logger.Info("app was gracefully shutdown")
}

func (a *app) initGateway(ctx context.Context) error {
	mux := runtime.NewServeMux()

	conn, err := grpc.DialContext(
		ctx,
		a.publicServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	if err := pb.RegisterFileServiceHandlerClient(
		ctx,
		mux,
		pb.NewFileServiceClient(conn),
	); err != nil {
		return err
	}

	a.gatewayServer = &http.Server{
		Addr:    a.gatewayAddress,
		Handler: mux,
	}

	return nil
}
