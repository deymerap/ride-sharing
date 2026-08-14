package main

import (
	pb "github.com/deymerap/ride-sharing/shared/proto/trip"
	"github.com/deymerap/ride-sharing/shared/types"
)

type PreviewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type TripStartRequest struct {
	TripID    string `json:"trip_id"`
	DriverID  string `json:"driver_id"`
	StartTime string `json:"start_time"`
}

type CreateTripRequest struct {
	RideFareID string `json:"rideFareID"`
	UserID     string `json:"userID"`
	StartTime  string `json:"startTime"`
}

func (p *PreviewTripRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  p.Destination.Latitude,
			Longitude: p.Destination.Longitude,
		},
	}
}

func (t *TripStartRequest) ToProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		TripID:    t.TripID,
		DriverID:  t.DriverID,
		StartTime: t.StartTime,
	}
}

func (c *CreateTripRequest) ToProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{}
}
