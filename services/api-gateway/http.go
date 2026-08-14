package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deymerap/ride-sharing/services/api-gateway/grpc_clients"
	"github.com/deymerap/ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody PreviewTripRequest
	if err := decodeJSONRequest(r, &reqBody); err != nil {
		http.Error(w, "failed to parse JSON request", http.StatusBadRequest)
		return
	}

	// validate the request body
	if reqBody.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	// GRPC call to trip-service
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal("failed to create trip service client:", err)
		http.Error(w, "failed to create trip service client", http.StatusInternalServerError)
	}
	defer tripService.Close()

	respBody, err := tripService.Client.PreviewTrip(r.Context(), reqBody.ToProto())
	if err != nil {
		log.Println("failed to call trip-service:", err)
		http.Error(w, "failed to call trip-service", http.StatusInternalServerError)
		return
	}
	response := contracts.APIResponse{Data: respBody}
	log.Printf("Response from trip-service : %+v", response)
	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	log.Println("RECEIVED REQUEST TO START TRIP IN API GATEWAY")

	var reqBody CreateTripRequest
	if err := decodeJSONRequest(r, &reqBody); err != nil {
		http.Error(w, "failed to parse JSON request", http.StatusBadRequest)
		return
	}

	log.Printf("REQUEST BODY IN API GATEWAY: %+v", reqBody)

	// validate the request body
	if reqBody.RideFareID == "" {
		http.Error(w, "rideFareID is required", http.StatusBadRequest)
		return
	}
	if reqBody.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	log.Println("CREATE GRPC CLIENT IN API GATEWAY")
	//gRPC call to trip-service
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal("failed to create trip service client:", err)
		http.Error(w, "failed to create trip service client", http.StatusInternalServerError)
	}
	defer tripService.Close()

	// Call the CreateTrip method on the trip service
	respBody, err := tripService.Client.CreateTrip(r.Context(), reqBody.ToProto())

	response := contracts.APIResponse{Data: respBody}
	log.Printf("Response from trip-service : %+v", response)
	writeJSON(w, http.StatusCreated, response)
}

func decodeJSONRequest(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}
