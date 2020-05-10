package greeting

import (
	"context"
	"fmt"

	pb "github.com/dgyoshi/grpc-ratelimiter-example/internal/pkg/pb/greeting"
	"google.golang.org/grpc"
)

type GreetingServiceClient struct {
	name   string
	client pb.GreetingClient
}

func NewClient(host, port, name string) *GreetingServiceClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure())

	if err != nil {
		panic(err.Error())
	}

	client := pb.NewGreetingClient(conn)

	return &GreetingServiceClient{
		name:   name,
		client: client,
	}
}

func (c *GreetingServiceClient) Hello(ctx context.Context, msg string) (string, error) {
	res, err := c.client.Hello(ctx, &pb.Msg{
		Name: c.name,
		Msg:  msg,
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return res.GetMsg(), nil

}
