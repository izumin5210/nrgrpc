# nrgrpc
[![Build Status](https://travis-ci.com/izumin5210/nrgrpc.svg?branch=master)](https://travis-ci.com/izumin5210/nrgrpc)
[![codecov](https://codecov.io/gh/izumin5210/nrgrpc/branch/master/graph/badge.svg)](https://codecov.io/gh/izumin5210/nrgrpc)
[![GoDoc](https://godoc.org/github.com/izumin5210/nrgrpc?status.svg)](https://godoc.org/github.com/izumin5210/nrgrpc)
[![Go project version](https://badge.fury.io/go/github.com%2Fizumin5210%2Fnrgrpc.svg)](https://badge.fury.io/go/github.com%2Fizumin5210%2Fnrgrpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/izumin5210/nrgrpc)](https://goreportcard.com/report/github.com/izumin5210/nrgrpc)
[![license](https://img.shields.io/github/license/izumin5210/nrgrpc.svg)](./LICENSE)

gRPC `stats.Handler` implementation to measure and send performances metrics to New Relic.

## Example
### gRPC server

`nrgrpc.NewServerStatsHandler` creates a [`stats.Handler`](https://godoc.org/google.golang.org/grpc/stats#Handler) instance for gRPC servers.
When the handler is passed to a gRPC server with [`stats.StatsHandler`](https://godoc.org/google.golang.org/grpc#StatsHandler),
it will set a [`newrelic.Transaction`](https://godoc.org/github.com/newrelic/go-agent#Transaction) into a request `context.Context` using [`nrutil.SetTransaction`](https://godoc.org/github.com/izumin5210/newrelic-contrib-go/nrutil#SetTransaction).
So you can retrieve `newrelic.Transaction` instances with [`nrutil.Transaction`](https://godoc.org/github.com/izumin5210/newrelic-contrib-go/nrutil#Transaction).

```go
func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	// Initiailze a `newrelic.Appliation`
	nrapp, err := newrelic.NewApplication(newrelic.Config{
		AppName: "your_app",
		License: "your_license_key",
	})
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		// Create a `stats.Handler` from `newrelic.Application`
		stats.StatsHandler(nrgrpc.NewServerStatsHandler(nrapp)),
	)

	// Register server implementations

	s.Serve(lis)
}
```

### gRPC client

```go
func main() {
	// Initiailze a `newrelic.Appliation`
	nrapp, err := newrelic.NewApplication(newrelic.Config{
		AppName: "your_app",
		License: "your_license_key",
	})
	if err != nil {
		panic(err)
	}

	// Create a `grpc.ClientConn` with `stats.Handler`
	conn, err := grpc.Dial(
		":3000",
		grpc.WithInsecure(),
		grpc.WithStatsHandler(nrgrpc.NewClientStatsHandler()),
	)
	if err != nil {
		panic(err)
	}


	// Register http handler using `github.com/izumin5210/newrelic-contrib-go/nrhttp`.
	// This wrapper sets `newrelic.Transaction` into the `http.Request`'s context.
	nrhttp.WrapHandleFunc(app, "/foo", func(w http.ResponseWriter, r *http.Request) {
		resp, err := NewFooServiceClient.BarCall(r.Context(), &BarRequest{})
		// ...
	})
}
```
