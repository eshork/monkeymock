package monkeymock

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type mockMethodCallInterface interface {
	Call(methodName string, args ...interface{}) []interface{}
}

type mockMethodContainerStruct struct {
	mockMethodSetupInterface

	lastmockMethodStructPtr *mockMethodStruct     // pointer to the last referenced mockMethod
	mockMethodPtrs          [](*mockMethodStruct) // ordered list of mockMethod declarations
}

type mockMethodStruct struct {
	parentMockStruct *mockStruct // reference back to the parent
	methodName       string

	// callable tracking
	callCountExpected     int                 // times this mockMethod was expected to be called (defaults to 0)
	expectedArgsValues    methodArgumentsList //
	expectedArgsAny       bool
	expectedReturnsValues methodReturnsList
	callRecords           [](*callRecordStruct)
	callOriginal          bool // when true, indicates the original method implementation should be called by the Mock
}

type methodArgumentsList []interface{}
type methodReturnsList []interface{}
type callRecordStruct struct {
	givenArgs       methodArgumentsList
	receivedReturns methodReturnsList
}

func (m *mockStruct) assertMethods(t *testing.T) {
	t.Helper()
	for _, mockMethodPtr := range m.mockMethodPtrs {
		mockMethodPtr.assertMethod(t)
	}
}

// times this mockMethod has been called
func (m *mockMethodStruct) calledCount() int {
	return len(m.callRecords)
}

func (m *mockMethodStruct) assertMethod(t *testing.T) {
	t.Helper()

	//
	// t.Errorf("\nyay itsa me!!!!")

	// number of calls
	m.assertMethodCallCount(t)
}

// support calls
// validate calls
// assert calls

func (m *mockMethodStruct) assertMethodCallCount(t *testing.T) {
	t.Helper()
	switch actualCalls := m.calledCount(); {
	case m.callCountExpected == 0: // was a 'Maybe' expectation
	case m.callCountExpected == -1 && actualCalls > 0: // failed 'Never' expectation
		methodName := stringifyMethodName(m)
		withargs := stringifyMethodArgs(m)
		withreturns := stringifyMethodReturns(m)
		tFail(t, fmt.Sprintf("Method called more than expected: \n"+
			"method  : %s\n"+
			"          withargs   : %s\n"+
			"          withreturns: %s\n"+
			"expected: Never\n"+
			"actual  : %d", methodName, withargs, withreturns, actualCalls))
	case actualCalls > m.callCountExpected: // called more than expected
		methodName := stringifyMethodName(m)
		withargs := stringifyMethodArgs(m)
		withreturns := stringifyMethodReturns(m)
		tFail(t, fmt.Sprintf("Method called more than expected: \n"+
			"method  : %s\n"+
			"          withargs   : %s\n"+
			"          withreturns: %s\n"+
			"expected: %d\n"+
			"actual  : %d", methodName, withargs, withreturns, m.callCountExpected, actualCalls))
	case actualCalls < m.callCountExpected: // called less than expected
		methodName := stringifyMethodName(m)
		withargs := stringifyMethodArgs(m)
		withreturns := stringifyMethodReturns(m)
		tFail(t, fmt.Sprintf("Method called less than expected: \n"+
			"method  : %s\n"+
			"          withargs   : %s\n"+
			"          withreturns: %s\n"+
			"expected: %d\n"+
			"actual  : %d", methodName, withargs, withreturns, m.callCountExpected, actualCalls))
	}
}

// Call mocked method instance with the given args.
// This call mechanism simulates a real call onto a Mock, and may induce
// a subsequent real call into the underlying object if required.
// It will trigger the expectations of the Mock, making it a useful for safe-typed mocks.
func (m *mockStruct) Call(methodName string, args ...interface{}) []interface{} {
	// find the referenced method -- this includes finging the most appropriate signature
	for _, mockMethodPtr := range m.mockMethodPtrs {
		if mockMethodPtr.methodName == methodName {
			return mockMethodPtr.call(args)
		}
	}
	panicMockMethodNotFound(getHumanTypeName(m.mockedObjectRef), methodName)
	return []interface{}{false} // satisfying pedantic compiler; this line is never reached...
}

// enact a call against a specific mockMethod
func (m *mockMethodStruct) call(args methodArgumentsList) methodReturnsList {
	var retVals methodReturnsList

	// record the call and given args
	callRecord := new(callRecordStruct)
	m.callRecords = append(m.callRecords, callRecord)
	callRecord.givenArgs = copyInterfaceList(args)

	// if this call was not expected (ie Never) then we should panic now
	m.panicIfCallExpectedNever()

	// validate args against method signature
	m.validateMethodCallArgsSignature(args) // TODO: fix this

	// if no declared return pattern and not AndCallsOriginal or AndCallsFunc, needs to panic now
	if !m.callOriginal {
		if len(m.expectedReturnsValues) == 0 { // no viable returns values!!!
			panicMockMethodReturnsNotDefined(stringifyMethodName(m))
		}
	}

	// should fall through to original function?
	if m.callOriginal {
		// more handy local variable ref
		objectRef := m.parentMockStruct.mockedObjectRef

		// get a usable method handle
		methodHandle := getObjectMethodByName(objectRef, m.methodName)

		// run the actual method call and try not to blow up
		retVals = callObjectMethodByName(methodHandle, objectRef, args)

		// capture the return values from the function
		callRecord.receivedReturns = copyInterfaceList(retVals)
	}

	// should fall through to custom handler function?
	//   monkeymock protect needed?

	// has expected return?
	//   yes - give expected, log actual
	//   no - give what we really received
	if len(m.expectedReturnsValues) > 0 {
		retVals = copyInterfaceList(m.expectedReturnsValues)
	}

	// return []interface{}{false, false}
	// return []interface{}{7}
	return retVals
}

//
func (m *mockMethodStruct) panicIfCallExpectedNever() {
	if m.callCountExpected == -1 {
		panicMockMethodCalledButNeverExpected(m)
	}
}

func (m *mockMethodStruct) reflectMethodType() reflect.Type {
	method, found := reflect.TypeOf(m.parentMockStruct.mockedObjectRef).MethodByName(m.methodName)
	if !found {
		panic("DERP!")
	}
	return method.Type
}

func (m *mockMethodStruct) validateMethodCallArgsSignature(args methodArgumentsList) {
	givenArgsSigStr := typeListToString(args)

	// generate argsSig from mockMethodStruct ptr
	var realArgSig []string
	{
		methodType := m.reflectMethodType()
		iMax := methodType.NumIn()
		for i := 0; iMax > i; i++ {
			inArg := methodType.In(i)
			realArgSig = append(realArgSig, inArg.String())
		}
	}
	// drop the first arg (it will always be a ref-ptr to the object's struct)
	realArgSig = realArgSig[1:]

	// generate a sig string
	realArgSigStr := ""
	for _, v := range realArgSig {
		realArgSigStr += "<" + v + ">, "
	}

	// panic on mismatch
	realArgSigStr = strings.Trim(realArgSigStr, " ,")
	if realArgSigStr != givenArgsSigStr {
		panicMockMethodCallInvalidArgsSignature(m, realArgSigStr, givenArgsSigStr)
	}
}
