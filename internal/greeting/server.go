package greeting

import (
	"context"
	"fmt"

	pb "github.com/dgyoshi/grpc-ratelimiter-example/internal/pkg/pb/greeting"
)

type Server struct {
	name string
}

func NewServer(name string) *Server {
	return &Server{
		name: name,
	}
}

func (s *Server) Hello(
	ctx context.Context,
	req *pb.Msg,
) (
	*pb.Reply,
	error,
) {
	return &pb.Reply{Msg: fmt.Sprintf("Hi %s, this is %s.", req.GetName(), s.name)}, nil
}
