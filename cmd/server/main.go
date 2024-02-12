package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/danblok/auth/client"
	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/logging"
	"github.com/danblok/auth/internal/service"
	"github.com/danblok/auth/proto"
)

func main() {
	httpAddr := flag.String("http", ":3000", "listen addr of the http server")
	grpcAddr := flag.String("grpc", ":4000", "listen addr of the grpc server")
	keyPath := flag.String("keypath", "/run/secrets/jwt_key", "key path of a signing jwt key")
	ctx := context.Background()
	flag.Parse()

	key, err := os.ReadFile(*keyPath)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewJWTService([]byte(key))
	svc = logging.NewLoggingService(svc)

	grpcClient, err := client.NewGRPCClient(*grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			tokenResp, err := grpcClient.Token(ctx, &proto.TokenRequest{Payload: "some payload"})
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%+v\n", tokenResp)

			validateResp, err := grpcClient.Validate(ctx, &proto.ValidateRequest{Token: "aoenuth"})
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%+v\n", validateResp)
		}
	}()

	go api.NewGRPCServer(svc).Serve(*grpcAddr)

	httpServer := api.NewHTTPServer(svc, *httpAddr)
	log.Printf("started server on http://localhost%s\n", *httpAddr)
	log.Printf(`available routes:
	receive token: POST http://localhost%s/token {"payload": "mypayload"}
	validate token: GET http://localhost%s/validate?token=<your_token>`, *httpAddr, *httpAddr)
	if err := httpServer.Run(); err != nil {
		log.Fatal("couldn't start the server")
	}
}
