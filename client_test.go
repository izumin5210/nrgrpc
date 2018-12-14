package nrgrpc

import (
	"context"
	"testing"

	nrgrpctesting "github.com/izumin5210/nrgrpc/testing"
	"github.com/newrelic/go-agent"
)

func Test_Client(t *testing.T) {
	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetClientStatsHandler(NewClientStatsHandler())
	ctx.Setup()
	defer ctx.Teardown()

	txn := nrgrpctesting.NewFakeNRTxn(t, "/echo")

	resp, err := ctx.Client.Echo(newrelic.NewContext(context.TODO(), txn), &nrgrpctesting.EchoRequest{Message: "foobar"})

	if resp == nil {
		t.Error("The request should return a response")
	}

	if err != nil {
		t.Errorf("should be no errors, but got %v", err)
	}

	if got, want := txn.Name, "/echo"; got != want {
		t.Errorf("the transaction name is %q, want %q", got, want)
	}
}

func Test_Gateway(t *testing.T) {
	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetClientStatsHandler(NewGatewayStatsHandler())
	ctx.Setup()
	defer ctx.Teardown()

	txn := nrgrpctesting.NewFakeNRTxn(t, "/echo")

	resp, err := ctx.Client.Echo(newrelic.NewContext(context.TODO(), txn), &nrgrpctesting.EchoRequest{Message: "foobar"})

	if resp == nil {
		t.Error("The request should return a response")
	}

	if err != nil {
		t.Errorf("should be no errors, but got %v", err)
	}

	if got, want := txn.Name, "/testing.TestService/Echo"; got != want {
		t.Errorf("the transaction name is %q, want %q", got, want)
	}
}
