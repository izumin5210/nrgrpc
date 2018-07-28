package nrgrpc

import (
	"testing"

	"github.com/izumin5210/newrelic-contrib-go/nrutil"
	nrgrpctesting "github.com/izumin5210/nrgrpc/testing"
	"golang.org/x/net/context"
)

func Test_Client(t *testing.T) {
	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetClientStatsHandler(NewClientStatsHandler())
	ctx.Setup()
	defer ctx.Teardown()

	txn := nrgrpctesting.NewFakeNRTxn(t, "/echo")

	resp, err := ctx.Client.Echo(nrutil.SetTransaction(context.TODO(), txn), &nrgrpctesting.EchoRequest{Message: "foobar"})

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

	resp, err := ctx.Client.Echo(nrutil.SetTransaction(context.TODO(), txn), &nrgrpctesting.EchoRequest{Message: "foobar"})

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
