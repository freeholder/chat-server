package main

import (
	"context"
	"log"
	"time"

	desc "github.com/freeholder/chat-server/pkg/note_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "172.26.112.94:50051"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: #{err}")
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Create(ctx, &desc.CreateRequest{Usernames: []string{"Ivan", "Petr"}})
	if err != nil {
		log.Fatalf("failed to create usernames list")
	}

	log.Printf("Users list created with ID: %d", r.GetId())
}
