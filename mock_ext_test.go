package monkeymock_test

// Tests that we can run from an external perspective

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/eshork/monkeymock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExampleExternalInterface interface {
	ExamplePublicMethod(trash1 string, trash2 int) int
}
type ExampleExternalStruct struct {
	ExampleExternalInterface
	someInt int
}

func (ex *ExampleExternalStruct) ExamplePublicMethod(trash1 string, trash2 int) int {
	if trash1 == "panic" {
		panic("i was told to panic!")
	}
	return trash2
}

func TestMockExpectationBuilder(t *testing.T) {
	suite.Run(t, new(testMockExpectationBuilder))
}

type testMockExpectationBuilder struct {
	suite.Suite
	fakeT *testing.T
}

func (s *testMockExpectationBuilder) SetupTest() {
	s.fakeT = new(testing.T)
}

func (s *testMockExpectationBuilder) TestExpectReturnsMockInterface() {
	exampleStruct := struct{}{}
	assert.Implements(s.T(), new(monkeymock.Mock), monkeymock.Expect(exampleStruct))
}

func (s *testMockExpectationBuilder) TestEmptyMockPassesAssert() {
	exampleStruct := struct{}{}
	require.False(s.T(), s.fakeT.Failed())
	mock := monkeymock.Expect(exampleStruct)
	mock.AssertExpectations(s.fakeT)
	assert.False(s.T(), s.fakeT.Failed())
}

func (s *testMockExpectationBuilder) TestUncalledMockOnceFailsAssert() {
	mock := monkeymock.Expect(ExampleExternalStruct{})
	mock.ToReceive("ExamplePublicMethod").Once()
	// mock.AssertExpectations(s.T())
	mock.AssertExpectations(s.fakeT)
	assert.True(s.T(), s.fakeT.Failed())
}

func (s *testMockExpectationBuilder) TestCalledMockMoreThanNeverPanics() {
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod").Never()
	assert.Panics(s.T(), func() {
		mock.Call("ExamplePublicMethod", "7", 8)
	})
}

func (s *testMockExpectationBuilder) TestCalledMockTooMuchFailsAssert() {
	// require.False(s.T(), s.fakeT.Failed())
	// mock := monkeymock.Expect(struct{}{})

	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj)
	mock.ToReceive("ExamplePublicMethod").Once().AndCallsOriginal()
	mock.Call("ExamplePublicMethod", "7", 8)
	mock.Call("ExamplePublicMethod", "7", 8)
	mock.AssertExpectations(s.fakeT)
	assert.True(s.T(), s.fakeT.Failed())
}

func (s *testMockExpectationBuilder) TestCalledMockOncePassesAssert() {
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj)
	mock.ToReceive("ExamplePublicMethod").Once().WithReturns(19)
	mock.Call("ExamplePublicMethod", "7", 8)
	mock.AssertExpectations(s.fakeT)
	assert.False(s.T(), s.fakeT.Failed())
}

func (s *testMockExpectationBuilder) TestUnmockableTypePanics() {
	assert.Panics(s.T(), func() {
		var _ = monkeymock.Expect(7)
	})
}

func (s *testMockExpectationBuilder) TestUnmockableTypePtrPanics() {
	assert.Panics(s.T(), func() {
		v := 7
		var _ = monkeymock.Expect(&v)
	})
}

func (s *testMockExpectationBuilder) TestMissingToReceiveMethodPanics() {
	exampleStruct := struct{}{}
	mock := monkeymock.Expect(exampleStruct)
	assert.Panics(s.T(), func() {
		mock.ToReceive("EmptyStructsDontHaveMethods")
	})
}

func (s *testMockExpectationBuilder) TestWithArgsRequiresToReceiveOrPanics() {
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj)
	assert.Panics(s.T(), func() {
		mock.WithArgs(7, 7)
	})
}

func (s *testMockExpectationBuilder) TestValidSetupWithArgsDoesNotPanic() {
	// valid calls to WithArgs() should not panic
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")

	assert.NotPanics(s.T(), func() {
		mock.WithArgs("taco", 7)
	})
}

func (s *testMockExpectationBuilder) TestInvalidSetupWithArgsDoesPanic() {
	// invalid calls to WithArgs() should  panic
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
	assert.Panics(s.T(), func() {
		mock.WithArgs(7, 7)
	}, "Calling WithArgs with invalid args patterns should panic")
}

func (s *testMockExpectationBuilder) TestWithArgsPanicsIfArgsAlreadyDeclared() {
	{ // WithArgs after WithArgs
		testObj := ExampleExternalInterface(&ExampleExternalStruct{})
		mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
		mock.WithArgs("taco", 7)
		assert.Panics(s.T(), func() {
			mock.WithArgs("taco", 7)
		}, "Calling WithArgs after WithArgs already set should panic")
	}
	{ // WithArgs after WithAnyArgs
		testObj := ExampleExternalInterface(&ExampleExternalStruct{})
		mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
		mock.WithAnyArgs()
		assert.Panics(s.T(), func() {
			mock.WithArgs("taco", 7)
		}, "Calling WithArgs after WithAnyArgs already set should panic")
	}
}

func (s *testMockExpectationBuilder) TestWithReturnsOverridesActual() {
	// check original
	{
		someNumber := gofakeit.Number(1, 99)
		testObj := ExampleExternalInterface(&ExampleExternalStruct{})
		mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod").
			WithArgs("anything", someNumber).
			AndCallsOriginal()
		retVal := mock.Call("ExamplePublicMethod", "anything", someNumber)
		assert.Equal(s.T(), []interface{}{someNumber}, retVal) // expect a single 6
	}
	// check override
	{
		someNumber := gofakeit.Number(1, 99)
		someOtherNumber := gofakeit.Number(100, 199)
		testObj := ExampleExternalInterface(&ExampleExternalStruct{})
		mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod").
			WithArgs("anything", someOtherNumber).
			WithReturns(someNumber).
			AndCallsOriginal()
		retVal := mock.Call("ExamplePublicMethod", "anything", someOtherNumber)
		assert.Equal(s.T(), []interface{}{someNumber}, retVal) // expect a single 99
	}
}

// Mock Doubles cannot work until named exported methods can be injected into runtime defined structs
func TestMockDoubles(t *testing.T) {
	// suite.Run(t, new(testMockDoubles))
}

type testMockDoubles struct {
	suite.Suite
	fakeT *testing.T
}

func (s *testMockDoubles) SetupTest() {
	s.fakeT = new(testing.T)
}

func (s *testMockDoubles) TestDoubleIsNotOriginalObject() {
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
	mockDbl := mock.AsDouble()
	assert.NotEqual(s.T(), testObj, mockDbl)
}

func (s *testMockDoubles) TestDoubleImplementsInterface() {
	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
	require.Implements(s.T(), (*ExampleExternalInterface)(nil), testObj)
	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
	mockDbl := mock.AsDouble()
	assert.Implements(s.T(), (*ExampleExternalInterface)(nil), mockDbl)
	var _ = mockDbl.(ExampleExternalInterface)
}

// func (s *testMockDoubles) xTestDoubleIsNotOriginalObject() {
// 	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
// 	require.Implements(s.T(), (*ExampleExternalInterface)(nil), testObj)

// 	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
// 	mockDbl := mock.AsDouble()
// 	assert.Implements(s.T(), (*ExampleExternalInterface)(nil), mockDbl)
// 	// AssertNotSame(s.T(), testObj, mockDbl)
// }

// func (s *testMockDoubles) xTestDoubleCanBeAssignedToInterface() {
// 	testObj := ExampleExternalInterface(&ExampleExternalStruct{})
// 	mock := monkeymock.Expect(testObj).ToReceive("ExamplePublicMethod")
// 	mockDbl := mock.AsDouble()
// 	// var testInterface ExampleExternalInterface
// 	testObj = mockDbl.(ExampleExternalInterface)
// 	// testInterface = mockDbl
// 	if false { // make compiler happy; we don't want to test actual calls just yet
// 		// testInterface.ExamplePublicMethod("just", 1)
// 	}
// }
