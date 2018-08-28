package nrgrpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"

	nrgrpctesting "github.com/izumin5210/nrgrpc/testing"
)

func Test_Server_RequestSuccessfully(t *testing.T) {
	app := nrgrpctesting.NewFakeNRApp(t)

	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetServerStatsHandler(NewServerStatsHandler(app, WithIgnoredMethods("/testing.TestService/Empty")))
	ctx.Setup()
	defer ctx.Teardown()

	resp, err := ctx.Client.Echo(context.TODO(), &nrgrpctesting.EchoRequest{Message: "foobar"})

	if resp == nil {
		t.Error("The request should return a response")
	}

	if err != nil {
		t.Errorf("should be no errors, but got %v", err)
	}

	app.CheckTxnCount(1)

	txn := app.Txns[0]

	txn.CheckEnded()
	txn.CheckStatusCode(codes.OK)
	txn.CheckIgnored(false)
	txn.CheckErrorCount(0)
}

func Test_Server_RequestFailed(t *testing.T) {
	app := nrgrpctesting.NewFakeNRApp(t)

	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetServerStatsHandler(NewServerStatsHandler(app, WithIgnoredMethods("/testing.TestService/Empty")))
	ctx.Setup()
	defer ctx.Teardown()

	resp, err := ctx.Client.Error(context.TODO(), new(empty.Empty))

	if resp != nil {
		t.Errorf("should not return a response, but got %v", resp)
	}

	if err == nil {
		t.Error("should return errors")
	}

	app.CheckTxnCount(1)

	txn := app.Txns[0]

	txn.CheckEnded()
	txn.CheckStatusCode(codes.Internal)
	txn.CheckIgnored(false)
	txn.CheckErrorCount(1)
}

func Test_Server_Ignored(t *testing.T) {
	app := nrgrpctesting.NewFakeNRApp(t)

	ctx := nrgrpctesting.CreateTestContext(t)
	ctx.SetServerStatsHandler(NewServerStatsHandler(app, WithIgnoredMethods("/testing.TestService/Empty")))
	ctx.Setup()
	defer ctx.Teardown()

	resp, err := ctx.Client.Empty(context.TODO(), new(empty.Empty))

	if resp == nil {
		t.Error("The request should return a response")
	}

	if err != nil {
		t.Errorf("should be no errors, but got %v", err)
	}

	app.CheckTxnCount(1)

	txn := app.Txns[0]

	txn.CheckEnded()
	txn.CheckStatusCode(codes.OK)
	txn.CheckIgnored(true)
	txn.CheckErrorCount(0)
}
