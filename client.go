package nrgrpc

import (
	"github.com/izumin5210/newrelic-contrib-go/nrutil"
	newrelic "github.com/newrelic/go-agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc/stats"
)

// NewClientStatsHandler creates a new stats.Handler instance for measuring application performances with New Relic.
func NewClientStatsHandler() stats.Handler {
	return &clientStatsHandlerImpl{}
}

type clientStatsHandlerImpl struct {
}

func (h *clientStatsHandlerImpl) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	txn := nrutil.Transaction(ctx)
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

type ctxKeyClientSegment struct{}

func setSegment(ctx context.Context, seg newrelic.Segment) context.Context {
	return context.WithValue(ctx, ctxKeyClientSegment{}, seg)
}

func getSegment(ctx context.Context) (st newrelic.Segment, ok bool) {
	if v := ctx.Value(ctxKeyClientSegment{}); v != nil {
		st, ok = v.(newrelic.Segment)
	}
	return
}
