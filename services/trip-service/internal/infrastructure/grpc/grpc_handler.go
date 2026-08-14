package grpc

import (
	"context"
	"log"

	"github.com/deymerap/ride-sharing/services/trip-service/internal/domain"
	pb "github.com/deymerap/ride-sharing/shared/proto/trip"
	"github.com/deymerap/ride-sharing/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	pickupCoord := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}
	destinationCoord := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}

	route, err := h.service.GetRoute(ctx, pickupCoord, destinationCoord)
	if err != nil {
		log.Printf("failed to get route: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	//1. Estimate the ride fares prices based on the route ex: distance
	//2. Store the ride fares for the create trip to fetch and validate.

	return &pb.PreviewTripResponse{
		Route: route.ToProto(),
		RideFares: []*pb.RideFares{
			{
				Id:                "ride_fare_1",
				UserID:            req.UserID,
				PackageSlug:       "sedan",
				TotalPriceInCents: 1000,
			},
			{
				Id:                "ride_fare_2",
				UserID:            req.UserID,
				PackageSlug:       "suv",
				TotalPriceInCents: 1500,
			},
		},
	}, nil
}

func (h *gRPCHandler) StartTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	// Implement the logic to start a trip here

	log.Printf("CREATE TRIP REQUEST IN GPRC HANDLER: %+v", req)

	return &pb.CreateTripResponse{
		TripID:    req.TripID,
		DriverID:  req.DriverID,
		StartTime: req.StartTime,
	}, nil
}
