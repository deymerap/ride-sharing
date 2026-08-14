package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deymerap/ride-sharing/shared/contracts"
	"github.com/deymerap/ride-sharing/shared/util"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer wsConn.Close() // close the WebSocket connection when the function returns

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID query parameter is required", http.StatusBadRequest)
		return
	}

	for {
		// Read message from WebSocket
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			http.Error(w, "Failed to read message from WebSocket", http.StatusInternalServerError)
			break
		}

		// Process the message (you can add your own logic here)
		//responseMessage := "Received: " + string(message)
		log.Printf("User %s sent message: %s", userID, message)

		// // Write response back to WebSocket
		// err = wsConn.WriteMessage(websocket.TextMessage, []byte(responseMessage))
		// if err != nil {
		// 	http.Error(w, "Failed to write message to WebSocket", http.StatusInternalServerError)
		// 	return
		// }
	}
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("STEP 1: Handling rider WebSocket connection...")
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer wsConn.Close() // close the WebSocket connection when the function returns
	log.Println("STEP 2: WebSocket connection established")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID query parameter is required", http.StatusBadRequest)
		return
	}
	log.Println("STEP 3: User ID is valid")

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		http.Error(w, "packageSlug query parameter is required", http.StatusBadRequest)
		return
	}
	log.Println("STEP 4: Package slug is valid")

	type Driver struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}
	log.Println("STEP 5: Driver struct initialized")

	msgDriver := Driver{
		ID:             userID,
		Name:           "Deymer Perea",
		ProfilePicture: util.GetRandomAvatar(1),
		CarPlate:       "ABC-787",
		PackageSlug:    packageSlug,
	}
	log.Println("STEP 6: Driver message initialized")
	dri, err := json.Marshal(msgDriver)
	if err != nil {
		panic(err)
	}
	log.Printf("CREATED DRIVER: %s", dri)

	msg := contracts.WSDriverMessage{
		Type: "driver.cmd.register",
		Data: dri,
	}

	log.Printf("CREATED DRIVER MESSAGE: %v", msg)

	if err := wsConn.WriteJSON(msg); err != nil {
		http.Error(w, "Failed to write JSON message to WebSocket", http.StatusInternalServerError)
		return
	}
	log.Println("STEP 7: Driver message sent to WebSocket")
	for {
		// Read message from WebSocket
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			http.Error(w, "Failed to read message from WebSocket", http.StatusInternalServerError)
			break
		}
		log.Printf("User %s sent message: %s", userID, message)
	}
}
