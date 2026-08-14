package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/deymerap/ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"github.com/deymerap/ride-sharing/services/trip-service/internal/infrastructure/repository"
	"github.com/deymerap/ride-sharing/services/trip-service/internal/service"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	//ctx := context.Background()
	inmemRepo := repository.NewInMemRepository()
	svc := service.NewService(inmemRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", GrpcAddr, err)
	}

	// iniciar el servidor gRPC
	grpcServer := grpcserver.NewServer()
	// TODO: inicializar la implementación del manejador gRPC
	grpc.NewGRPCHandler(grpcServer, svc)

	log.Printf("Starting gRPC server on %s", lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
			cancel()
		}
	}()

	// Esperar a que se reciba la señal de cancelación
	<-ctx.Done()
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}
