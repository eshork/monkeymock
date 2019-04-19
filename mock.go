package monkeymock

// https://medium.com/@utter_babbage/breaking-the-type-system-in-golang-aka-dynamic-types-8b86c35d897b

// Mock represents a single mock instance configuration.
// The Mock configuration object itself is not usable as a stand-in for future
// function calls, but several methods are available to assist with those needs.
// Ref: [Mock.Double()], Mock.Partial()
type Mock interface {
	mockAssertions
	mockMethodSetupInterface
	mockMethodCallInterface
	mockDoubleInterface
	mockPartialInterface
	// mockCallableInterface
	// mockCallCounterInterface
}

// internal struct implementing the Mock interface
type mockStruct struct {
	Mock
	mockMethodContainerStruct

	mockedObjectRef interface{}
}

// Expect is the first step to building an expectation around a thing, either a type or an object
// instance. It will return a Mock object, which can then be used to define the details of the
// expectation.
// Each call to Expect creates a new Mock and thus begins defining a new expectation that can be evaluated
// for completion accuracy en masse (the typical method) via  ([AssertExpectations]), or on an individual
// basis by directly calling ([Mock.Assert]) for each mock you'd like to evaluate.
func Expect(refObject interface{}) Mock {
	validateIsMockableObjectRef(refObject)
	mock := new(mockStruct)          // every Expect is a new assert condition...
	mock.mockedObjectRef = refObject // store a reference to the original object/interface
	appendToMockList(mock)           // throw it onto the FIFO stack...
	return mock                      // make condition stacking easy...
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

//// Global Mock list management - as every Expect is supposed to generate an expectation by default!

//

var gTheMockList []Mock

// adds the given Mock to the end of the current MockList (starts a new list if necessary)
func appendToMockList(newmock Mock) {
	gTheMockList = append(gTheMockList, newmock)
}

// removes the given Mock from the MockList (may result in 0 size list; you don't care)
func removeFromMockList(mock Mock) {
	if gTheMockList == nil {
		return
	}
	for i, v := range gTheMockList {
		if mock == v {
			if i == 0 {
				gTheMockList = gTheMockList[1:] // trim the first item off the list
				return
			}
			if lastI := len(gTheMockList) - 1; i == lastI {
				gTheMockList = gTheMockList[:lastI] // trim the last item off the list
				return
			}
			newList := make([]Mock, len(gTheMockList)-1)
			copy(newList, gTheMockList[:i])
			copy(newList[i:], gTheMockList[i+1:])
			gTheMockList = newList
		}
	}
}

// drops the current MockList entirely, so a new one may be started
func clearMockList() {
	gTheMockList = nil // that was easy...
}

func sizeOfMockList() int {
	return len(gTheMockList)
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// will panic if not a suitable object type
func validateIsMockableObjectRef(refObject interface{}) {
	mockable, typename := isMockableObjectRef(refObject)
	if !mockable {
		panicUnmockableType(typename)
	}
}
