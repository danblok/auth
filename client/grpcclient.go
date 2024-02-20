package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/danblok/auth/proto"
)

// NewGRPCClient returns GRPC client to communicate with the TokenService GRPC server.
func NewGRPCClient(addr string) (proto.TokenServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return proto.NewTokenServiceClient(conn), nil
}

// NewGRPCClientTLS returns GRPC client to communicate with the TokenService GRPC server securely.
func NewGRPCClientTLS(addr string, creds credentials.TransportCredentials) (proto.TokenServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return proto.NewTokenServiceClient(conn), nil
}
