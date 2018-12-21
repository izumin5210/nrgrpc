package nrgrpc

import (
	"context"
	"net/http"
	"net/url"

	newrelic "github.com/newrelic/go-agent"
	"google.golang.org/grpc/metadata"
)

type request struct {
	url    *url.URL
	header http.Header
}

var (
	mappedHeaders = []struct{ from, to string }{
		{from: newrelic.DistributedTracePayloadHeader},
		{from: "user-agent"},
		{from: "x-request-start", to: "x-request-start"},
		{from: "x-queue-start", to: "x-queue-start"},
		{from: "grpcgateway-x-request-start", to: "x-request-start"},
		{from: "grpcgateway-x-queue-start", to: "x-queue-start"},
	}
)

func newRequest(ctx context.Context, fullMethodName string) newrelic.WebRequest {
	h := http.Header{}
	h.Add("content-type", "application/grpc")
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for _, m := range mappedHeaders {
			from := m.from
			to := m.to
			if to == "" {
				to = from
			}
			if v := md.Get(from); len(v) > 0 {
				h.Add(to, v[0])
			}
		}
	}

	return &request{
		url:    &url.URL{Path: fullMethodName},
		header: h,
	}
}

func (r *request) Header() http.Header {
	return r.header
}

func (r *request) URL() *url.URL {
	return r.url
}

func (r *request) Method() string {
	return ""
}

func (r *request) Transport() newrelic.TransportType {
	// TODO: should define `TransportGRPC"
	return newrelic.TransportUnknown
}
