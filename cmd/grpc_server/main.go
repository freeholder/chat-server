package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	desc "github.com/freeholder/chat-server/pkg/note_v1"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

type server struct {
	desc.UnimplementedNoteV1Server
	db *pgx.Conn
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {

	var id int64

	err := s.db.QueryRow(
		ctx,
		"INSERT INTO request (usernames) VALUES ($1) RETURNING id",
		req.Usernames,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to insert", err)
	}

	// res, err := s.db.Exec(ctx, "INSERT INTO request (usernames) VALUES ($1)", req.Usernames)
	// if err != nil {
	// 	log.Fatalf("failed to insert note: %v", err)
	// }

	// log.Printf("inserted %d rows", res.RowsAffected())

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_, err := s.db.Exec(
		ctx,
		"DELETE FROM request WHERE id = $1",
		req.Id,
	)

	if err != nil {
		fmt.Printf("Delete error - %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func main() {

	ctx := context.Background()
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	fmt.Println(color.GreenString("Server has started!"))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{db: con})

	log.Printf("gRPC server listening on port %s", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// res, err := con.Exec(ctx, "INSERT INTO request (usernames) VALUES ($1)", []string{"Egor", "Maxim"})
	// if err != nil {
	// 	log.Fatalf("failed to insert note: %v", err)
	// }

	// log.Printf("inserted %d rows", res.RowsAffected())

}
