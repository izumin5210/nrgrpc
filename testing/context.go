package testing

import (
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// TestContext is a testing helper for generating a gRPC server and a gRPC client.
type TestContext struct {
	t *testing.T

	ServerOpts []grpc.ServerOption
	ClientOpts []grpc.DialOption

	serverListener net.Listener
	server         *grpc.Server
	clientConn     *grpc.ClientConn

	Service TestServiceServer
	Client  TestServiceClient
}

// CreateTestContext returns a new TestContext object.
func CreateTestContext(t *testing.T) *TestContext {
	return &TestContext{
		t:       t,
		Service: &testServiceServer{},
	}
}

// SetServerStatsHandler sets stats.Handler to a test server.
func (c *TestContext) SetServerStatsHandler(h stats.Handler) {
	c.ServerOpts = append(c.ServerOpts, grpc.StatsHandler(h))
}

// SetClientStatsHandler sets stats.Handler to a test client.
func (c *TestContext) SetClientStatsHandler(h stats.Handler) {
	c.ClientOpts = append(c.ClientOpts, grpc.WithStatsHandler(h))
}

// Setup starts a server and creates a client connection.
func (c *TestContext) Setup() {
	if c.Service == nil {
		c.t.Fatal("Should set errorstesting.TestService implementaiton")
	}
	c.setupServer()
	c.setupClient()
}

// Teardown disconnects a client connection and stops a server
func (c *TestContext) Teardown() {
	time.Sleep(10 * time.Millisecond)
	if c.serverListener != nil {
		c.server.GracefulStop()
		c.serverListener.Close()
	}
	if c.clientConn != nil {
		c.clientConn.Close()
	}
}

func (c *TestContext) setupServer() {
	var err error
	c.serverListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		c.Teardown()
		c.t.Fatal("Failed to listen local network")
	}
	c.server = grpc.NewServer(c.ServerOpts...)
	RegisterTestServiceServer(c.server, c.Service)
	go c.server.Serve(c.serverListener)
}

func (c *TestContext) setupClient() {
	var err error
	dialOpts := append(c.ClientOpts, grpc.WithBlock(), grpc.WithTimeout(2*time.Second), grpc.WithInsecure())
	c.clientConn, err = grpc.Dial(c.serverListener.Addr().String(), dialOpts...)
	if err != nil {
		c.Teardown()
		c.t.Fatal("Failed to create a client connection")
	}
	c.Client = NewTestServiceClient(c.clientConn)
}
