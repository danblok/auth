package api

import (
	"context"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/danblok/auth/pkg/types"
	"github.com/danblok/auth/proto"
)

type GRPCTokenServer struct {
	proto.UnimplementedTokenServiceServer
	svc types.TokenService
}

func NewGRPCServer(svc types.TokenService) *GRPCTokenServer {
	return &GRPCTokenServer{
		svc: svc,
	}
}

func (s *GRPCTokenServer) Serve(addr string) error {
	tokenServiceServer := NewGRPCServer(s.svc)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterTokenServiceServer(grpcServer, tokenServiceServer)

	return grpcServer.Serve(ln)
}

func (s *GRPCTokenServer) Token(ctx context.Context, req *proto.TokenRequest) (*proto.TokenResponse, error) {
	reqId := uuid.NewString()
	ctx = context.WithValue(ctx, types.RequestID("request_id"), reqId)
	token, err := s.svc.Token(ctx, []byte(req.Payload))
	if err != nil {
		return nil, err
	}

	return &proto.TokenResponse{Token: string(token)}, nil
}

func (s *GRPCTokenServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	reqId := uuid.NewString()
	ctx = context.WithValue(ctx, types.RequestID("request_id"), reqId)
	err := s.svc.Validate(ctx, []byte(req.Token))
	if err != nil {
		return &proto.ValidateResponse{Valid: false}, nil
	}

	return &proto.ValidateResponse{Valid: true}, nil
}
