package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/deymerap/ride-sharing/services/trip-service/internal/domain"
	localTypes "github.com/deymerap/ride-sharing/services/trip-service/pkg/types"
	"github.com/deymerap/ride-sharing/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	log.Println("CREATING TRIP IN TRIP SERVICE")
	trip := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	return s.repo.CreateTrip(ctx, trip)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*localTypes.OSRMResponse, error) {
	baseURL := "http://router.project-osrm.org/route/v1/driving/"
	url := fmt.Sprintf("%s%f,%f;%f,%f?overview=full&geometries=geojson", baseURL,
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get route from OSRM: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OSRM response body: %w", err)
	}
	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("OSRM returned non-200 status code: %d", resp.StatusCode)
	// }

	var osrmResponse localTypes.OSRMResponse
	if err := json.Unmarshal(body, &osrmResponse); err != nil {
		return nil, fmt.Errorf("failed to decode OSRM response: %w", err)
	}

	if len(osrmResponse.Routes) == 0 {
		return nil, fmt.Errorf("no routes found in OSRM response")
	}

	return &osrmResponse, nil
}
