package monkeymock

import (
	"fmt"
)

func panicRuntimeGeneral(msg string) {
	panicMsg := fmt.Sprintf("\n"+
		"GENERAL PANIC: %s\n",
		msg)
	tPanicMockRuntime(panicMsg)
}

func panicUnexpectedMethodCall(m *mockMethodStruct) {
	panicMsg := fmt.Sprintf("\n" +
		"Unexpected call to Method (no handler defined for Mock Double): \n" +
		"object    : \n" +
		"method    : \n" +
		"args found: \n")
	tPanicMockRuntime(panicMsg)
}

func panicExpectationDeclaredBeforeToReceive(srcMethod string) {
	panicMsg := fmt.Sprintf("\n"+
		"mock.%s called before mock.ToReceive()\n"+
		"ToReceive method must be declared before other expectations\n",
		srcMethod)
	tPanicMockSetup(panicMsg)
}

func panicArgsAlreadyDeclared(srcMethod string) {
	panicMsg := fmt.Sprintf("\n"+
		"Cannot call mock.%s within a mock.ToReceive() with already declared arguments.\n",
		srcMethod)
	tPanicMockSetup(panicMsg)
}

func panicReturnsAlreadyDeclared(srcMethod string) {
	panicMsg := fmt.Sprintf("\n"+
		"Cannot call mock.%s within a mock.ToReceive() with already declared returns.\n",
		srcMethod)
	tPanicMockSetup(panicMsg)
}

func panicUnmockableType(typeName string) {
	panicMsg := fmt.Sprintf("\n"+
		"Expect attempted on an unmockable object type: \n"+
		"found type: <%s>\n",
		typeName)
	tPanicMockSetup(panicMsg)
}

func panicMethodNotFoundInObjectRef(typeName string, methodName string) {
	panicMsg := fmt.Sprintf("\n"+
		"Cannot create expectation: ToReceive(\"%s\")\n"+
		"Method not found within refObject type: <%s>\n",
		methodName, typeName)
	tPanicMockSetup(panicMsg)
}

func panicMockMethodNotFound(typeName string, methodName string) {
	panicMsg := fmt.Sprintf("\n"+
		"Method not found within Mock: <%s>.%s\n",
		typeName, methodName)
	tPanicMockRuntime(panicMsg)
}

func panicMockMethodReturnsNotDefined(fullMethodName string) {
	panicMsg := fmt.Sprintf("\n"+
		"Method called without return value declaration within Mock: %s\n"+
		"Must either specify a return value (WithReturns) or fallback\n"+
		"to original method implementation (AndCallsOriginal)\n",
		fullMethodName)
	tPanicMockRuntime(panicMsg)
}

func panicMockMethodCallInvalidArgsSignature(mockMethod *mockMethodStruct, expectedSig string, receivedSig string) {
	methodName := stringifyMethodName(mockMethod)
	tPanicMockRuntime(fmt.Sprintf("Mock Method called with invalid arguments signature: \n"+
		"method  : %s\n"+
		"types expected : %s\n"+
		"types received : %s\n"+
		"",
		methodName, expectedSig, receivedSig))
}

func panicMockMethodCalledButNeverExpected(m *mockMethodStruct) {
	methodName := stringifyMethodName(m)
	withargs := stringifyMethodArgs(m)
	withreturns := stringifyMethodReturns(m)
	tPanicMockRuntime(fmt.Sprintf("Method called but Never was expected: \n"+
		"method  : %s\n"+
		"          withargs   : %s\n"+
		"          withreturns: %s\n",
		methodName, withargs, withreturns))
}

func panicMockWithArgsMismatch(mockMethod *mockMethodStruct, expectedArgs string, receivedArgs string) {
	methodName := stringifyMethodName(mockMethod)
	tPanicMockSetup(fmt.Sprintf("mock.WithArgs() called with mismatched signature: \n"+
		"method  : %s\n"+
		"types expected : %s\n"+
		"types received : %s\n"+
		"",
		methodName, expectedArgs, receivedArgs))
}

func panicAsDoubleNotImplemented() {
	panicMsg := fmt.Sprintf("\n" +
		"Mock Doubles \"AsDouble()\" cannot currently be implemented.\n" +
		"Additional support is necessary within the language.\n" +
		"ref: https://github.com/golang/go/issues/16522 \n")
	tPanicMockSetup(panicMsg)
}
