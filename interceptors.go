package nrgrpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/izumin5210/newrelic-contrib-go/nrutil"
	"github.com/newrelic/go-agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type key struct{}

var (
	txnKey key
)

// UnaryServerInterceptor returns a new unary server interceptor to set newrelic transaction
func UnaryServerInterceptor(app newrelic.Application, optFuncs ...Option) grpc.UnaryServerInterceptor {
	opts := composeOptions(optFuncs)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if opts.IsIgnored(info.FullMethod) {
			return handler(ctx, req)
		}
		txn := app.StartTransaction(info.FullMethod, nil, nil)
		defer txn.End()
		resp, err := handler(nrutil.SetTransaction(ctx, txn), req)
		if err != nil {
			txn.NoticeError(err)
		}
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor to set newrelic transaction
func StreamServerInterceptor(app newrelic.Application, optFuncs ...Option) grpc.StreamServerInterceptor {
	opts := composeOptions(optFuncs)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if opts.IsIgnored(info.FullMethod) {
			return handler(srv, stream)
		}
		txn := app.StartTransaction(info.FullMethod, nil, nil)
		defer txn.End()
		wrappedStream := grpc_middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = nrutil.SetTransaction(wrappedStream.Context(), txn)
		err := handler(srv, wrappedStream)
		if err != nil {
			txn.NoticeError(err)
		}
		return err
	}
}
