package main

import (
	"context"
	"fmt"

	soipb "github.com/koooyooo/soi-go/pkg/srv/server/grpc"
	"google.golang.org/grpc"
)

func main() {
	if err := control(); err != nil {
		fmt.Println(err)
	}
}

func control() error {
	fmt.Println("Hello gRPC")
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := soipb.NewSoiServiceClient(conn)
	ctx := context.Background()

	req := soipb.SoiRequest{
		Data: &soipb.SoiData{
			Name:  "Name",
			Title: "Title",
			Uri:   "URL",
		},
	}

	for i := 0; i < 100; i++ {
		resp, err := cli.RegisterSoi(ctx, &req)
		if err != nil {
			return err
		}
		fmt.Print(resp.Status, " ")
	}
	return nil
}
