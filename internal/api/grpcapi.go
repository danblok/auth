package api

import (
	"context"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/danblok/auth/pkg/types"
	"github.com/danblok/auth/proto"
)

// GRPCTokenServer implements TokenService via GRPC transport.
type GRPCTokenServer struct {
	proto.UnimplementedTokenServiceServer
	svc types.TokenService
}

// NewGRPCServer creates new GRPC server.
func NewGRPCServer(svc types.TokenService) *GRPCTokenServer {
	return &GRPCTokenServer{
		svc: svc,
	}
}

// Serve runs grpc server.
func (s *GRPCTokenServer) Serve(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterTokenServiceServer(grpcServer, s)

	return grpcServer.Serve(ln)
}

// Token provides API on behalf of the GRPC server to receive token.
func (s *GRPCTokenServer) Token(ctx context.Context, req *proto.TokenRequest) (*proto.TokenResponse, error) {
	reqID := uuid.NewString()
	ctx = context.WithValue(ctx, types.RequestID("request_id"), reqID)
	token, err := s.svc.Token(ctx, []byte(req.Payload))
	if err != nil {
		return nil, err
	}

	return &proto.TokenResponse{Token: string(token)}, nil
}

// Validate provides API on behalf of the GRPC server to validate token.
func (s *GRPCTokenServer) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	reqID := uuid.NewString()
	ctx = context.WithValue(ctx, types.RequestID("request_id"), reqID)
	err := s.svc.Validate(ctx, []byte(req.Token))
	if err != nil {
		return &proto.ValidateResponse{Valid: false}, nil
	}

	return &proto.ValidateResponse{Valid: true}, nil
}
