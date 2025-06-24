package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	api "github.com/Open-Myth/protoc-gen-http-client/pb/greeting"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	mux := runtime.NewServeMux()

	ctx := context.Background()
	service := greetingService{}
	api.RegisterGreetingServiceHandlerServer(ctx, mux, &service)
	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mux)
}

type greetingService struct {
	api.UnimplementedGreetingServiceServer
}

func (*greetingService) SayHello(ctx context.Context, req *api.SayHelloRequest) (*api.SayHelloResponse, error) {
	fmt.Println(req.Ping)
	return &api.SayHelloResponse{
		Pong: "Hello " + req.Ping,
	}, nil
}
