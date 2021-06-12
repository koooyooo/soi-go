package main

import (
	"context"
	"fmt"
	"log"
	"net"

	soipb "github.com/koooyooo/soi-go/pkg/srv/server/grpc"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Run Server")

	listener, err := net.Listen("tcp", "0.0.0.0:5051")
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	soipb.RegisterSoiServiceServer(srv, &soiServer{})

	srv.Serve(listener)
}

type soiServer struct {
	soipb.UnimplementedSoiServiceServer
}

func (s *soiServer) RegisterSoi(c context.Context, r *soipb.SoiRequest) (*soipb.SoiResponse, error) {
	return &soipb.SoiResponse{
		Status: 1,
	}, nil
}
