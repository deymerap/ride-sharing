package grpc_clients

import (
	"os"

	pb "github.com/deymerap/ride-sharing/shared/proto/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tripServiceClient struct {
	Client     pb.TripServiceClient
	Connection *grpc.ClientConn
}

func NewTripServiceClient() (*tripServiceClient, error) {
	tripServicesURL := os.Getenv("TRIP_SERVICE_URL") // se usa el nombre del servicio en lugar de la IP para que funcione con docker-compose
	if tripServicesURL == "" {
		tripServicesURL = "trip-service:9093" // valor por defecto si no se encuentra la variable de entorno
	}
	conn, err := grpc.NewClient(tripServicesURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &tripServiceClient{
		Client:     pb.NewTripServiceClient(conn),
		Connection: conn}, nil
}

func (c *tripServiceClient) Close() {
	if c.Connection != nil {
		if err := c.Connection.Close(); err != nil {
			return
		}
	}
}
