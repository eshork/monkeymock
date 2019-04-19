package monkeymock

import (
	"fmt"
	"testing"
)

func tFail(t *testing.T, failureMessage string) {
	t.Helper()
	t.Errorf("MonkeyMock.AssertExpectations failed\n%s\n", failureMessage)
}
func tFatal(t *testing.T, failureMessage string) {
	t.Helper()
	t.Fatalf("MonkeyMock.AssertExpectations failed (FATAL)\n%s\n", failureMessage)
}
func tPanicMockRuntime(failureMessage string) {
	panicMsg := fmt.Sprintf("MonkeyMock PANIC at runtime: \n"+
		"%s\n", failureMessage)
	panic(panicMsg)
}
func tPanicMockSetup(failureMessage string) {
	panicMsg := fmt.Sprintf("MonkeyMock PANIC during Mock setup: \n"+
		"%s\n", failureMessage)
	panic(panicMsg)
}
