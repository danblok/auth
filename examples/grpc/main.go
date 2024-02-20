package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/danblok/auth/client"
	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/logging"
	"github.com/danblok/auth/internal/service"
	"github.com/danblok/auth/proto"
)

var (
	grpcAddr       = flag.String("grpc", ":3000", "Listen addr of the http server")
	jwtKeyPath     = flag.String("jwtkey", "data/jwt", "Key path of a signing jwt key")
	caCertPath     = flag.String("cacert", "data/ca.crt", "CA certificate path")
	serverCertPath = flag.String("srvcert", "data/server.crt", "Server certificate path")
	serverKeyPath  = flag.String("srvkey", "data/server.key", "Server private key path")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	key, err := os.ReadFile(*jwtKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewJWTService([]byte(key))
	svc = logging.NewLoggingService(svc)

	go func() {
		cert, err := tls.LoadX509KeyPair(*serverCertPath, *serverKeyPath)
		if err != nil {
			log.Fatalf("couldn't load x509 key pair: %v", err)
		}
		log.Printf("started GRPC server on [::]%s", *grpcAddr)
		if err := api.NewGRPCServer(svc).ServeTLS(*grpcAddr, cert); err != nil {
			log.Fatalf("couldn't run GRPC server: %v", err)
		}
	}()

	creds, err := credentials.NewClientTLSFromFile(*caCertPath, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := client.NewGRPCClientTLS(*grpcAddr, creds)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(2 * time.Second)
		tokenResp, err := client.Token(ctx, &proto.TokenRequest{Payload: "some payload"})
		if err != nil {
			log.Fatalf("GRPC client.Token: %v", err)
		}
		log.Printf("GRPC: %+v\n", tokenResp)

		validateResp, err := client.Validate(ctx, &proto.ValidateRequest{Token: tokenResp.Token})
		if err != nil {
			log.Fatalf("GRPC client.Validate: %v", err)
		}
		log.Printf("GRPC: %+v\n", validateResp)
	}
}
