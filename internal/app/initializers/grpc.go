package initializers

import (
	"fmt"
	"github.com/SiriusServiceDesk/gateway-service/pkg/notification_v1"
	"github.com/SiriusServiceDesk/notification-service/internal/app/dependencies"
	"github.com/SiriusServiceDesk/notification-service/internal/config"
	"github.com/SiriusServiceDesk/notification-service/internal/grpc/server"
	"github.com/SiriusServiceDesk/notification-service/internal/grpc/server/handlers"
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

func StartGRPCServer(grpcServer *server.GRPCServer) {
	go func() {
		err := grpcServer.Start()
		if err != nil {
			logger.Fatal("failed starting GRPC server", "error", err)
		}
	}()
}

func InitializeGRPCListener() (net.Listener, error) {
	cfg := config.GetConfig().GrpcServer
	logger.Info(fmt.Sprintf("GRPC starting on %s:%s", cfg.Host, cfg.Port))
	listener, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
	}
	return listener, nil
}

func InitializeGRPCServer(listener net.Listener, container *dependencies.Container) *server.GRPCServer {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	notification_v1.RegisterNotificationV1Server(
		grpcServer,
		handlers.NewHandler(container.MessageService, container.Queue))

	logger.Info("Complete register handlers")
	return server.NewGRPCServer(listener, grpcServer)
}
