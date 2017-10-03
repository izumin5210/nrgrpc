package nrgrpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/newrelic/go-agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type key struct{}

var (
	txnKey key
)

// UnaryServerInterceptor returns a new unary server interceptor to set newrelic transaction
func UnaryServerInterceptor(app newrelic.Application) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		txn := app.StartTransaction(info.FullMethod, nil, nil)
		defer txn.End()
		return handler(setTransaction(ctx, txn), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor to set newrelic transaction
func StreamServerInterceptor(app newrelic.Application) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		txn := app.StartTransaction(info.FullMethod, nil, nil)
		defer txn.End()
		wrappedStream := grpc_middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = setTransaction(wrappedStream.Context(), txn)
		return handler(srv, stream)
	}
}

// Transaction extracts newrelic transaction object from request context
func Transaction(ctx context.Context) newrelic.Transaction {
	v := ctx.Value(txnKey)
	if v == nil {
		return nil
	}
	return v.(newrelic.Transaction)
}

func setTransaction(ctx context.Context, txn newrelic.Transaction) context.Context {
	return context.WithValue(ctx, txnKey, txn)
}
