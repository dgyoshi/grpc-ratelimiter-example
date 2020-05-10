package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/dgyoshi/grpc-ratelimiter-example/internal/greeting"
	"google.golang.org/grpc/metadata"
)

type config struct {
	ServerHost string `env:"SERVER_HOST,required"`
	ServerPort string `env:"SERVER_PORT,required"`
	Name       string `env:"CLIENT_NAME,required"`
	Interval   string `env:"CALL_INTERVAL,required"`
	Times      int    `env:"CALL_TIMES,required"`
}

func main() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		panic(err)
	}

	client := greeting.NewClient(conf.ServerHost, conf.ServerPort, conf.Name)

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	md := metadata.Pairs("name", conf.Name)
	ctx = metadata.NewOutgoingContext(ctx, md)

	count := 0
	wg := sync.WaitGroup{}
	for {
		if count >= conf.Times {
			break
		}
		wg.Add(1)
		count++

		go func(ctx context.Context, name string, i int) {
			defer wg.Done()
			fmt.Printf("call '%d'\n", i)

			res, err := client.Hello(ctx, fmt.Sprintf("This is %s.", name))
			if err != nil {
				fmt.Printf("error: count '%d', %s\n", i, err.Error())
			}

			fmt.Printf("respons: count '%d', msg: '%s'\n", i, res)
		}(ctx, conf.Name, count)

		time.Sleep(interval)
	}

	wg.Wait()

	fmt.Printf("%s stopped greeting.", conf.Name)
}
