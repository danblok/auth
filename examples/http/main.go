package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"time"

	"github.com/danblok/auth/client"
	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/logging"
	"github.com/danblok/auth/internal/service"
)

var (
	httpAddr       = flag.String("http", ":3000", "Listen addr of the http server")
	keyPath        = flag.String("jwtkey", "data/jwt", "Key path of a signing jwt key")
	caCertPath     = flag.String("cacert", "data/ca.crt", "CA certificate path")
	serverCertPath = flag.String("srvcert", "data/server.crt", "Server certificate path")
	serverKeyPath  = flag.String("srvkey", "data/server.key", "Server private key path")
)

func main() {
	ctx := context.Background()
	key, err := os.ReadFile(*keyPath)
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
		httpServer, err := api.NewHTTPServerTLS(svc, *httpAddr, cert)
		if err != nil {
			log.Fatalf("couldn't create a new HTTP server: %v", err)
		}
		log.Printf("started HTTP server on [::]%s\n", *httpAddr)
		if err := httpServer.Run(); err != nil {
			log.Fatalf("couldn't run HTTP server: %v", err)
		}
	}()

	// Read CA certificate.
	cert, err := os.ReadFile(*caCertPath)
	if err != nil {
		log.Fatalf("couldn't read CA cert: %v", err)
	}

	// Creating new HTTP TLS client
	client, err := client.NewHTPPClientTLS("localhost"+*httpAddr, cert)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Every 2 seconds fetch a new token a validate it.
		time.Sleep(2 * time.Second)
		tokenResp, err := client.Token(ctx, []byte("some payload"))
		if err != nil {
			log.Fatalf("HTTP client.Token: %v", err)
		}

		log.Printf("HTTP: %+v\n", tokenResp)

		validateResp, err := client.Validate(ctx, []byte(tokenResp.Token))
		if err != nil {
			log.Fatalf("HTTP client.Validate: %v", err)
		}
		log.Printf("HTTP: %+v\n", validateResp)
	}
}
