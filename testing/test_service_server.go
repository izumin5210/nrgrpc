package testing

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testServiceServer struct {
}

func (s *testServiceServer) Echo(_ context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Message: req.GetMessage()}, nil
}

func (s *testServiceServer) Empty(context.Context, *empty.Empty) (*empty.Empty, error) {
	return new(empty.Empty), nil
}

func (s *testServiceServer) Error(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Error(codes.Internal, "an error occurred")
}
