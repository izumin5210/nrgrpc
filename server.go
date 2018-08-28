package nrgrpc

import (
	"context"

	"github.com/izumin5210/newrelic-contrib-go/nrutil"
	"github.com/newrelic/go-agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// NewServerStatsHandler creates a new stats.Handler instance for measuring application performances with New Relic.
func NewServerStatsHandler(app newrelic.Application, opts ...Option) stats.Handler {
	return &serverStatsHandlerImpl{
		app:  app,
		opts: composeOptions(opts),
	}
}

type serverStatsHandlerImpl struct {
	app  newrelic.Application
	opts Options
}

func (h *serverStatsHandlerImpl) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	txn := h.app.StartTransaction(info.FullMethodName, nil, nil)

	if h.opts.IsIgnored(info.FullMethodName) {
		txn.Ignore()
	}

	ctx = nrutil.SetTransaction(ctx, txn)

	return ctx
}

func (h *serverStatsHandlerImpl) HandleRPC(ctx context.Context, s stats.RPCStats) {
	switch s := s.(type) {
	case *stats.End:
		txn := nrutil.Transaction(ctx)
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			txn.AddAttribute("metadata", md)
		}
		if err := s.Error; err != nil {
			if st, ok := status.FromError(s.Error); ok {
				txn.AddAttribute("grpcStatusCode", st.Code())
			}
			txn.NoticeError(err)
		} else {
			txn.AddAttribute("grpcStatusCode", codes.OK)
		}
		txn.End()
	}
}

func (h *serverStatsHandlerImpl) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

func (h *serverStatsHandlerImpl) HandleConn(ctx context.Context, s stats.ConnStats) {
	// no-op
}
