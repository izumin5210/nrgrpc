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

func newRequest(ctx context.Context, fullMethodName string) newrelic.WebRequest {
	h := http.Header{}
	h.Add("content-type", "application/grpc")
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if v := md.Get(newrelic.DistributedTracePayloadHeader); len(v) > 0 {
			h.Add(newrelic.DistributedTracePayloadHeader, v[0])
		}
		if v := md.Get("user-agent"); len(v) > 0 {
			h.Add("user-agent", v[0])
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
