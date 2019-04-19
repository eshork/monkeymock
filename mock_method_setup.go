package monkeymock

type mockMethodSetupInterface interface {
	// expectation indicators regarding number of times a method should be called
	ToReceive(methodName string) Mock
	Once() Mock           // alias for Times(1)
	Twice() Mock          // alias for Times(2)
	Times(count int) Mock // specifies a hard counter for expected number of calls
	Maybe() Mock          // resets call expectation to unspecified state (zero or more times)
	Never() Mock          // expects the expectation that the method will never be called
	// Calls() int           // returns the number of times the method was called upon
	// Are we missing flexible call counters? Ie MoreTimesThan and LessTimesThan ??? (maybe addressable by the Calls() counter)

	WithArgs(args ...interface{}) Mock // match the method expectation with a particular set of argument values
	WithAnyArgs() Mock                 // match the method expectation regardless of argument values (ie, declare a default matcher)

	WithReturns(returnValues ...interface{}) Mock // expect particular return value(s); will override actual return values if also "AndCallsOriginal", but such a case also throws a failure during AssertExpections if the values do not align

	AndCallsOriginal() Mock // expectation will actually perform a call to the original implementaion
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// ToReceive sets up an expectation for a method call
func (m *mockStruct) ToReceive(methodName string) Mock {
	if m.lastmockMethodStructPtr != nil {
		panic("ToReceive called multiple times (this usage is NYI)")
	}

	// validate reference object/type supports method name, or throw a panic
	if !objectRefHasMethod(m.mockedObjectRef, methodName) {
		panicMethodNotFoundInObjectRef(getHumanTypeName(m.mockedObjectRef), methodName)
	}

	// set up the new method record
	newmockMethod := new(mockMethodStruct)
	newmockMethod.parentMockStruct = m
	newmockMethod.methodName = methodName
	newmockMethod.callCountExpected = 0
	m.mockMethodPtrs = append(m.mockMethodPtrs, newmockMethod)
	m.lastmockMethodStructPtr = newmockMethod
	return m
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// Once - expect the method once
func (m *mockStruct) Once() Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("Once()")
	}
	return m.Times(1)
}

// Twice - expect the method twice
func (m *mockStruct) Twice() Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("Twice()")
	}
	return m.Times(2)
}

// Times - expect the method count times
func (m *mockStruct) Times(count int) Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("Times()")
	}
	m.lastmockMethodStructPtr.callCountExpected = count
	return m
}

// Maybe - expect the method zero or more times
func (m *mockStruct) Maybe() Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("Maybe()")
	}
	return m.Times(0)
}

// Never - expect the method to never be called
func (m *mockStruct) Never() Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("Never()")
	}
	m.lastmockMethodStructPtr.callCountExpected = -1
	return m
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// WithArgs - expect the method to be called with a particular set of argument values
// When multiple ToReceive expectations exist, argument expectations provide a matching filter
// to line return expectations up with the given arguments.
// Note: The given expected argument values are copied by value (shallow); changing the underlying
// values during runtime may result in unexpected behaviour
func (m *mockStruct) WithArgs(args ...interface{}) Mock {
	// panic when ToReceive is missing
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("WithArgs()")
	}

	// panic when args expectation already set
	if m.lastmockMethodStructPtr.expectedArgsValues != nil || m.lastmockMethodStructPtr.expectedArgsAny {
		panicArgsAlreadyDeclared("WithArgs()")
	}

	// type check the args list -- will throw panic if they mismatch
	m.lastmockMethodStructPtr.ensureMethodArgs(args)

	// store a copy of the args list for later reference
	m.lastmockMethodStructPtr.expectedArgsValues = copyInterfaceList(args)

	return m
}

// func argsfromVariadic(variadic ...interface{})

// WithAnyArgs - sets an explicit non-expectation regarding arguments.
// If no known arguments pattern match the current call, the WithAnyArgs pattern will
// be used as the default handler.
// Because it matches all argument patterns, only one WithAnyArgs may currently be declared.
// Attempts to set multiple WithAnyArgs currently results in a setup-time panic.
func (m *mockStruct) WithAnyArgs() Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("WithAnyArgs()")
	}
	// panic when args expectation already set
	if m.lastmockMethodStructPtr.expectedArgsValues != nil || m.lastmockMethodStructPtr.expectedArgsAny {
		panicArgsAlreadyDeclared("WithArgs()")
	}
	m.lastmockMethodStructPtr.expectedArgsAny = true
	return m
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// WithReturns - sets an expectation of a specific return value (or values).
// If the mock includes `AndCallsOriginal()`, the original method will be called,
// but the value returned will be replaced with this given expectation. The mismatch
// will be surfaceable via a call to AssertExpectations
func (m *mockStruct) WithReturns(returnValues ...interface{}) Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("WithReturns()")
	}
	// panic when args expectation already set
	if m.lastmockMethodStructPtr.expectedReturnsValues != nil {
		panicReturnsAlreadyDeclared("WithReturns()")
	}

	// TODO: validate returns signature

	// store a copy of the returns list for later reference
	m.lastmockMethodStructPtr.expectedReturnsValues = copyInterfaceList(returnValues)

	return m
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// AndCallsOriginal - sets up the mock expectation to call the original
// method implementation when the mock is called. Can be combined with WithReturns()
// to validate the original implementation is producing the expected results, or
// can be used without WithReturns() to let the produced return values pass through
// untouched.
func (m *mockStruct) AndCallsOriginal() Mock {
	if m.lastmockMethodStructPtr == nil {
		panicExpectationDeclaredBeforeToReceive("AndCallsOriginal()")
	}

	// panic when args expectation already set
	if m.lastmockMethodStructPtr.callOriginal == true {
		panicReturnsAlreadyDeclared("WithReturns()")
	}

	m.lastmockMethodStructPtr.callOriginal = true
	return m
}

// AndCallsFunc -
