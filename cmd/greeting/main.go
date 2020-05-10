package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/dgyoshi/grpc-ratelimiter-example/internal/greeting"
	ratelimit_interceptor "github.com/dgyoshi/grpc-ratelimiter-example/internal/interceptor/ratelimiter"
	pb "github.com/dgyoshi/grpc-ratelimiter-example/internal/pkg/pb/greeting"
	"github.com/dgyoshi/grpc-ratelimiter-example/internal/pkg/pb/greeting/ratelimiter"
	"google.golang.org/grpc"
)

type config struct {
	Port     string `env:"SERVER_PORT,required"`
	Name     string `env:"SERVER_NAME,required"`
	Interval string `env:"REQUEST_INTERVAL,required"`
	Timeout  string `env:"REQUEST_TIMEOUT,required"`
	Capacity int64  `env:"REQUEST_CAPACITY,required"`
	Quantum  int64  `env:"REQUEST_QUANTUM,required"`
}

func main() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		panic(err)
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		panic(err)
	}
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		panic(err)
	}
	rateLimiter := ratelimiter.New(interval, timeout, conf.Capacity, conf.Quantum)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(ratelimit_interceptor.UnaryServerInterceptor(rateLimiter)))

	gs := greeting.NewServer(conf.Name)
	pb.RegisterGreetingServer(s, gs)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		panic(err.Error())
	}

	go func() {
		fmt.Printf("%s is listening on %s port\n", conf.Name, conf.Port)
		if err := s.Serve(listen); err != nil {
			panic(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for s := range c {
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return

		case syscall.SIGHUP:
		default:
			return
		}
	}

	return
}
