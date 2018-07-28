package nrgrpc

import (
	newrelic "github.com/newrelic/go-agent"
	"golang.org/x/net/context"
)

type ctxKeyClientSegment struct{}

func setSegment(ctx context.Context, seg newrelic.ExternalSegment) context.Context {
	return context.WithValue(ctx, ctxKeyClientSegment{}, seg)
}

func getSegment(ctx context.Context) (st newrelic.ExternalSegment, ok bool) {
	if v := ctx.Value(ctxKeyClientSegment{}); v != nil {
		st, ok = v.(newrelic.ExternalSegment)
	}
	return
}
