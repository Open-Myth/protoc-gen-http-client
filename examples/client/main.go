package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	api "github.com/Open-Myth/protoc-gen-http-client/pb/greeting"
)

func main() {
	for {
		client := api.NewGreetingServiceHTTPClient("http://localhost:8080")
		resp, err := client.SayHello(context.Background(), &api.SayHelloRequest{
			Ping: fmt.Sprintf("Hello %d", rand.Int31n(100)),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(err)
		fmt.Println(resp)
		time.Sleep(2 * time.Second)
	}
}
