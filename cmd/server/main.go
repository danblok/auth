package main

import (
	"flag"
	"log"
	"os"

	"github.com/danblok/auth/internal/api"
	"github.com/danblok/auth/internal/logging"
	"github.com/danblok/auth/internal/service"
)

func main() {
	addr := flag.String("addr", ":3000", "listen addr of the server")
	keyPath := flag.String("keypath", "/run/secrets/jwt_key", "key path of a signing jwt key")
	flag.Parse()

	key, err := os.ReadFile(*keyPath)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewJWTService([]byte(key))
	svc = logging.NewLoggingService(svc)

	srv := api.NewHTTPService(svc, *addr)
	log.Printf("started server on http://localhost%s\n", *addr)
	log.Printf(`available routes:
	receive token: POST http://localhost%s/token
	validate token: GET http://localhost%s/validate?token=<your_token>`, *addr, *addr)
	if err := srv.Run(); err != nil {
		log.Fatal("couldn't start the server")
	}
}
