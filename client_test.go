package nrgrpc

import (
	"testing"

	"github.com/izumin5210/newrelic-contrib-go/nrutil"
	nrgrpctesting "github.com/izumin5210/nrgrpc/testing"
	newrelic "github.com/newrelic/go-agent"
	"golang.org/x/net/context"
)

func Test_Client_RequestSuccessfully(t *testing.T) {
	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetClientStatsHandler(NewClientStatsHandler())
	ctx.Setup()
	defer ctx.Teardown()

	app, err := newrelic.NewApplication(newrelic.Config{AppName: "nrgrpctesting", Enabled: false})
	if err != nil {
		t.Fatalf("failed to create newrelic.Application: %v", err)
	}
	txn := app.StartTransaction("testendpoint", nil, nil)

	resp, err := ctx.Client.Echo(nrutil.SetTransaction(context.TODO(), txn), &nrgrpctesting.EchoRequest{Message: "foobar"})

	if resp == nil {
		t.Error("The request should return a response")
	}

	if err != nil {
		t.Errorf("should be no errors, but got %v", err)
	}
}
