package nrgrpc

import (
	"context"

	"github.com/izumin5210/newrelic-contrib-go/nrutil"
	"github.com/newrelic/go-agent"
	"google.golang.org/grpc/stats"
)

// NewClientStatsHandler creates a new stats.Handler instance for measuring application performances with New Relic.
func NewClientStatsHandler() stats.Handler {
	return &clientStatsHandlerImpl{}
}

// NewGatewayStatsHandler creates a new stats.Handler instance for measuring application performances with New Relic.
func NewGatewayStatsHandler() stats.Handler {
	return &clientStatsHandlerImpl{updateTxnName: true}
}

type clientStatsHandlerImpl struct {
	updateTxnName bool
}

func (h *clientStatsHandlerImpl) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	txn := nrutil.Transaction(ctx)

	if h.updateTxnName {
		txn.SetName(info.FullMethodName)
	}

	seg := newrelic.StartSegment(txn, info.FullMethodName)

	return setSegment(ctx, seg)
}

func (h *clientStatsHandlerImpl) HandleRPC(ctx context.Context, s stats.RPCStats) {
	switch s.(type) {
	case *stats.End:
		if seg, ok := getSegment(ctx); ok {
			seg.End()
		}
	}
}

func (h *clientStatsHandlerImpl) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

func (h *clientStatsHandlerImpl) HandleConn(ctx context.Context, s stats.ConnStats) {
	// no-op
}
