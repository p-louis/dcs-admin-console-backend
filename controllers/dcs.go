package controllers

import (
	"context"
	"flag"
	"github.com/DCS-gRPC/go-bindings/dcs/v0/mission"
	"github.com/DCS-gRPC/go-bindings/dcs/v0/net"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"time"
)

var playersCache = []*net.GetPlayersResponse_GetPlayerInfo{}

func refreshPlayersCache(netClient net.NetServiceClient) error {
	response, err := netClient.GetPlayers(context.Background(), &net.GetPlayersRequest{})
	if err != nil {
		return err
	}
	playersCache = response.Players
	return nil
}

func getPlayerName(id uint32) string {
	for _, player := range playersCache {
		if player.Id == id {
			return player.Name
		}
	}

	return ""
}

type MessageInput struct {
	Message string `json:"message" binding:"required"`
}

func SendChatMessage(c *gin.Context) {
	var input MessageInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var addr = flag.String("server", "localhost:50051", "address of the target DCS-gRPC server endpoint")
	var message = flag.String("message", input.Message, "chat message to send")

	flag.Parse()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithTimeout(5 * time.Second),
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Panicf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Printf("Sending message to chat...")
	netClient := net.NewNetServiceClient(conn)
	_, err = netClient.SendChat(context.Background(), &net.SendChatRequest{
		Message: *message,
	})
	if err != nil {
		log.Panicf("Failed to send message: %v", err)
	}

	missionClient := mission.NewMissionServiceClient(conn)
	stream, err := missionClient.StreamEvents(context.Background(), &mission.StreamEventsRequest{})
	if err != nil {
		log.Panicf("Failed to open events stream: %v", err)
	}

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Server closed events stream")
			return
		} else if err != nil {
			log.Panicf("Error reading from events stream: %v", err)
		}

		switch inner := event.GetEvent().(type) {
		case *mission.StreamEventsResponse_PlayerSendChat:
			playerName := getPlayerName(inner.PlayerSendChat.PlayerId)
			if playerName == "" {
				err := refreshPlayersCache(netClient)
				if err != nil {
					log.Panicf("Failed to refresh players list: %v", err)
				}
				playerName = getPlayerName(inner.PlayerSendChat.PlayerId)
			}

			log.Printf("[%v] %v", playerName, inner.PlayerSendChat.Message)
		}
	}

}
