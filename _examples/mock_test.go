package examples

// Example tests

import (
	"testing"

	// "github.com/eshork/monkeymock"
	"github.com/stretchr/testify/assert"
)

type yourStruct struct{}

func (s *yourStruct) yourmethod(input string) string {
	return "Original Return"
}

func TestOriginalObject(t *testing.T) {
	yourObj := new(yourStruct)
	assert.Equal(t, yourObj.yourmethod("in"), "Original Return")
}

// func xTestBasicMock(t *testing.T) {

// 	type yourStructMock struct {
// 		monkeymock.Mock
// 		yourStruct
// 	}
// 	yourObj := new(yourStructMock)

// 	// set up a mock atop an existing object

// 	mock := monkeymock.Expect(yourObj).(*yourStruct) //.
// 	// ToReceive("method").    // object method name (string symbol value)
// 	// Once().                 // how many times will this be called? once (obviously)
// 	// WithArgs("string_arg"). // expect specific arg values
// 	// AndCallOriginal().      // preserve the original functionality
// 	// WithReturns(7)          // expect a specific return value

// 	// do some test actions using the mock
// 	yourObj.yourmethod("somevalue")

// 	// assert that all established expectations have been met
// 	// monkeymock.AssertExpectations(t)

// 	// this specific example will fail with an args-based assertion failure!
// }
