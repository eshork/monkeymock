package monkeymock

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	. "github.com/eshork/monkeymock/testsupports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

type ExamplePartialInternalInterface interface {
	ExamplePublicMethod(trash1 string, trash2 int) int
}
type ExamplePartialInternalStruct struct {
	ExamplePartialInternalInterface
	someInt int
}

func (m *ExamplePartialInternalStruct) ExamplePublicMethod(trash1 string, trash2 int) int {
	return trash2
}
func (m ExamplePartialInternalStruct) ExamplePublicMethod2(trash1 string, trash2 int) int {
	return trash2
}
func (m *ExamplePartialInternalStruct) ExamplePublicMethod3(trash1 string, trash2 int) int {
	return trash2
}

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

func TestMockPartialInternals(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())
	suite.Run(t, new(testMockPartialInternals))
}

type testMockPartialInternals struct {
	suite.Suite
	fakeT *testing.T
}

func (s *testMockPartialInternals) AfterTest(_, _ string) {
	clearPartialObjectMethodIntercepts() // always clear out existing intercepts
	ClearExpectations(s.T())             // reset the expectations
}

func (s *testMockPartialInternals) TestCreatesPartialFromGivenObject() {
	exampleStruct := &ExamplePartialInternalStruct{}

	// try to make the partial
	partialObj := Expect(exampleStruct).AsPartial()

	// assert they are the same object
	AssertSame(s.T(), exampleStruct, partialObj)
}

func (s *testMockPartialInternals) TestPartialCreatesObjectMethodIntercept() {
	require.Zero(s.T(), len(interceptRecords), "interceptRecords should start empty")
	exampleStruct := &ExamplePartialInternalStruct{}
	var _ = Expect(exampleStruct).ToReceive("ExamplePublicMethod").AsPartial()
	assert.NotZero(s.T(), len(interceptRecords), "interceptRecords should have entry after partial created")
}

func (s *testMockPartialInternals) TestObjectMethodIntercepts() {
	exampleStruct := &ExamplePartialInternalStruct{}
	untouchedStruct := &ExamplePartialInternalStruct{}

	{ // both objects currently operate as normal
		n := gofakeit.Number(0, 99)
		assert.Equal(s.T(),
			untouchedStruct.ExamplePublicMethod("junk", n),
			exampleStruct.ExamplePublicMethod("junk", n),
			"Object methods should return equal values prior to mocking")
	}

	{ // mocked partial changes the output to the new expected value
		nOrig := gofakeit.Number(0, 99)
		nExpected := gofakeit.Number(100, 199)
		var _ = Expect(exampleStruct).
			ToReceive("ExamplePublicMethod").
			WithReturns(nExpected).
			AsPartial()
		ret := exampleStruct.ExamplePublicMethod("junk", nOrig)
		assert.Equal(s.T(), nExpected, ret,
			"Object method should return overridden value")
	}

	{ // both objects currently operate as normal
		nExpected := gofakeit.Number(200, 299)
		assert.Equal(s.T(), nExpected,
			untouchedStruct.ExamplePublicMethod("junk", nExpected),
			"Untouched Object methods should continue to execute per normal")
	}

	{ // methods without undeclared exceptions continue to operate as normal
		nExpected := gofakeit.Number(300, 399)
		assert.Equal(s.T(), nExpected,
			exampleStruct.ExamplePublicMethod3("junk", nExpected),
			"Uncaptured methods on AsPartial objects should continue to execute per normal")
	}
	// var _ = Expect(exampleStruct).
	// 	ToReceive("ExamplePublicMethod2").
	// 	AndCallsOriginal().
	// 	AsPartial()
	// exampleStruct.ExamplePublicMethod2("junk", 1)
	// assert they are the same object
	// AssertSame(s.T(), exampleStruct, partialObj)

}
