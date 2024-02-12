package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/danblok/auth/proto"
)

func NewGRPCClient(addr string) (proto.TokenServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return proto.NewTokenServiceClient(conn), nil
}
