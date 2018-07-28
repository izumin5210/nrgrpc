package testing

import (
	"net/http"
	"testing"

	newrelic "github.com/newrelic/go-agent"
	"google.golang.org/grpc/codes"
)

// NewFakeNRApp returns a new fake newrelic.Application instnace.
func NewFakeNRApp(t *testing.T) *FakeNRApp {
	return &FakeNRApp{t: t}
}

// FakeNRApp is an impplementation of newrelic.Application for tests.
type FakeNRApp struct {
	newrelic.Application
	Txns []*FakeNRTxn
	t    *testing.T
}

// StartTransaction implements the newrelic.Application interface.
func (a *FakeNRApp) StartTransaction(name string, _ http.ResponseWriter, _ *http.Request) newrelic.Transaction {
	txn := &FakeNRTxn{
		Name:  name,
		Attrs: make(map[string]interface{}),
		t:     a.t,
	}
	a.Txns = append(a.Txns, txn)
	return txn
}

// CheckTxnCount fails the test if created transaction count is invalid.
func (a *FakeNRApp) CheckTxnCount(c int) {
	a.t.Helper()
	if got, want := len(a.Txns), c; got != want {
		a.t.Fatalf("The application has %d transactions, want %d", got, want)
	}
}

// FakeNRTxn is an implementation of newrelic.Transaction for tests.
type FakeNRTxn struct {
	newrelic.Transaction
	Name    string
	Errs    []error
	Ended   bool
	Ignored bool
	Attrs   map[string]interface{}
	t       *testing.T
}

// End implements the newrelic.Transaction interface.
func (t *FakeNRTxn) End() error {
	t.Ended = true
	return nil
}

// Ignore implements the newrelic.Transaction interface.
func (t *FakeNRTxn) Ignore() error {
	t.Ignored = true
	return nil
}

// NoticeError implements the newrelic.Transaction interface.
func (t *FakeNRTxn) NoticeError(err error) error {
	t.Errs = append(t.Errs, err)
	return err
}

// AddAttribute implements the newrelic.Transaction interface.
func (t *FakeNRTxn) AddAttribute(k string, v interface{}) error {
	t.Attrs[k] = v
	return nil
}

// CheckEnded fails the test if the transaction has been ended.
func (t *FakeNRTxn) CheckEnded() {
	t.t.Helper()
	if !t.Ended {
		t.t.Error("the transaction should be ended")
	}
}

// CheckIgnored fails the test if the transaction's ignored status is invaild.
func (t *FakeNRTxn) CheckIgnored(ignored bool) {
	t.t.Helper()
	if got, want := t.Ignored, ignored; got != want {
		t.t.Errorf("the transaction was called Ignore %t, want %t", got, want)
	}
}

// CheckStatusCode fails the test if the transaction does not has a valid status code.
func (t *FakeNRTxn) CheckStatusCode(c codes.Code) {
	t.t.Helper()
	if code, ok := t.Attrs["grpcStatusCode"]; !ok {
		t.t.Error("attr grpcStatusCode is missing")
	} else if got, want := code, c; got != want {
		t.t.Errorf("grpcStatusCode is %v, want %v", got, want)
	}
}

// CheckErrorCount fails the test if the transaction's error  count is invalid.
func (t *FakeNRTxn) CheckErrorCount(c int) {
	t.t.Helper()
	if got, want := len(t.Errs), c; got != want {
		t.t.Fatalf("The transaction has %d errors, want %d", got, want)
	}
}
