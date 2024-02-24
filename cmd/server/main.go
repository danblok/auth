package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/logging"
	"github.com/danblok/auth/internal/service"
)

var (
	httpAddr       = flag.String("http", ":3000", "Listen addr of the http server")
	grpcAddr       = flag.String("grpc", ":4000", "Listen addr of the grpc server")
	jwtKeyPath     = flag.String("jwtkey", "/run/secrets/jwt_key", "Key path of a signing jwt key")
	serverCertPath = flag.String("srvcert", "/run/secrets/server_cert", "Server certificate path")
	serverKeyPath  = flag.String("srvkey", "/run/secrets/server_key", "Server private key path")
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	key, err := os.ReadFile(*jwtKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewJWTService([]byte(key))
	svc = logging.NewLoggingService(svc)

	eg := new(errgroup.Group)

	eg.Go(func() error {
		cert, err := tls.LoadX509KeyPair(*serverCertPath, *serverKeyPath)
		if err != nil {
			return err
		}
		log.Printf("started GRPC server on [::]%s", *grpcAddr)
		return api.NewGRPCServer(svc).ServeTLS(*grpcAddr, cert)
	})

	eg.Go(func() error {
		cert, err := tls.LoadX509KeyPair(*serverCertPath, *serverKeyPath)
		if err != nil {
			return err
		}
		httpServer, err := api.NewHTTPServerTLS(svc, *httpAddr, cert)
		if err != nil {
			return fmt.Errorf("couldn't create a new HTTP server: %v", err)
		}
		log.Printf("started HTTP server on [::]%s\n", *httpAddr)
		log.Printf(`available routes:
	receive token: POST [::]%s/token {"payload": "mypayload"}
	validate token: GET [::]%s/validate?token=<your_token>`, *httpAddr, *httpAddr)
		return httpServer.Run()
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
