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
	key, err := os.ReadFile("jwt")
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewJWTService([]byte(key))
	svc = logging.NewLoggingService(svc)
	addr := flag.String("addr", ":3000", "listen addr of the server")

	flag.Parse()

	srv := api.NewHTTPService(svc, *addr)
	log.Printf("started server on port: %s", *addr)
	if err := srv.Run(); err != nil {
		log.Fatal("couldn't start the server")
	}
}
