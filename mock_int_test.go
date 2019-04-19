package monkeymock

// Tests that we can only run from an internal perspective

import (

	// 	"errors"
	// 	"fmt"
	// 	"regexp"
	// 	"runtime"
	// 	"sync"

	"reflect"
	"testing"

	// 	"time"

	. "github.com/eshork/monkeymock/testsupports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// import (
// 	// 	"errors"
// 	// 	"fmt"
// 	// 	"regexp"
// 	// 	"runtime"
// 	// 	"sync"

// 	"testing"

// 	// 	"time"
// 	"github.com/eshork/monkeymock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// )

// each Expect increases global theMockList
// clearMockList resets theMockList

// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"runtime"
// 	"sync"

// 	"time"

// every Expect creates a unique Mock

type testInternalMocks struct {
	suite.Suite
	// fakeT *testing.T
}

func TestInternalMocks(t *testing.T) {
	// suite.Run(t, new(testInternalMocks))
}

func (s *testInternalMocks) TestMocksAreUniqueOnDiffRefObjects() {
	// two unique looking objects should be treated as unique
	mock1 := Expect(struct{}{}) // mock of an empty anonymous struct
	mock2 := Expect(struct{}{}) // mock of an empty anonymous struct
	assert.Implements(s.T(), new(Mock), mock1)
	assert.Implements(s.T(), new(Mock), mock2)
	AssertNotSame(s.T(), mock1, mock2)
}

func (s *testInternalMocks) TestHumanTypeNames() {
	assert.Equal(s.T(), getHumanTypeName(7), "int")
	i := 7
	assert.Equal(s.T(), getHumanTypeName(&i), "*int")

	assert.Equal(s.T(), getHumanTypeName(struct{}{}), "struct")
	assert.Equal(s.T(), getHumanTypeName(&struct{}{}), "*struct")

	type testStruct struct{}
	assert.Equal(s.T(), getHumanTypeName(testStruct{}), "testStruct")
	assert.Equal(s.T(), getHumanTypeName(&testStruct{}), "*testStruct")
}

type someExampleInterface interface {
	ExampleMethod(trash1 string, trash2 int) int
}
type someExampleStruct struct {
	someExampleInterface
}

func (s *someExampleStruct) ExampleMethod(trash1 string, trash2 int) int { return 8 }
func (s *someExampleStruct) ExampleMethod2()                             {}

func (s *testInternalMocks) TestGetMethodByName() {
	{
		testObj := someExampleInterface(&someExampleStruct{})
		methodHdl := getObjectMethodByName(testObj, "ExampleMethod")
		require.NotNil(s.T(), methodHdl)
		assert.Nil(s.T(), getObjectMethodByName(testObj, "ExampleMethod_missing"))
		ret := callObjectMethodByName(methodHdl, testObj, methodArgumentsList{"junk", 7})
		assert.EqualValues(s.T(), ret, methodReturnsList{8})
	}
	{
		testObj := &someExampleStruct{}
		methodHdl := getObjectMethodByName(testObj, "ExampleMethod")
		require.NotNil(s.T(), methodHdl)
		assert.Nil(s.T(), getObjectMethodByName(testObj, "ExampleMethod_missing"))
		ret := callObjectMethodByName(methodHdl, testObj, methodArgumentsList{"junk", 7})
		assert.EqualValues(s.T(), ret, methodReturnsList{8})
	}
}

func (s *testInternalMocks) TestVoidReturnsOnCalls() {
	{
		testObj := someExampleInterface(&someExampleStruct{})
		methodHdl := getObjectMethodByName(testObj, "ExampleMethod2")
		require.NotNil(s.T(), methodHdl)
		// both of these call formats are okay
		ret := callObjectMethodByName(methodHdl, testObj, methodArgumentsList{})
		assert.EqualValues(s.T(), ret, methodReturnsList{})
		ret = callObjectMethodByName(methodHdl, testObj, nil)
		assert.EqualValues(s.T(), ret, methodReturnsList{})
	}
}

func (s *testInternalMocks) TestGetNormalizedTypes() {
	type junkType struct {
		i int
		s string
	}

	{ // verify pointer input type resolves
		simpleObj := junkType{}
		simpleObjPtr := &simpleObj
		ptrType := reflect.TypeOf(simpleObjPtr)
		outPtrType, outObjType := getNormalizedObjectTypes(ptrType)

		require.NotNil(s.T(), outPtrType)
		require.NotNil(s.T(), outObjType)

		assert.Equal(s.T(), reflect.Ptr, outPtrType.Kind())
		assert.Equal(s.T(), reflect.Struct, outObjType.Kind())
	}
	{ // verify concrete input type resolves
		simpleObj := junkType{}
		objType := reflect.TypeOf(simpleObj)
		outPtrType, outObjType := getNormalizedObjectTypes(objType)

		require.NotNil(s.T(), outPtrType)
		require.NotNil(s.T(), outObjType)

		assert.Equal(s.T(), reflect.Ptr, outPtrType.Kind())
		assert.Equal(s.T(), reflect.Struct, outObjType.Kind())
	}

	{ // verify concrete input type resolves (non-struct)
		simpleObj := 7
		objType := reflect.TypeOf(simpleObj)
		outPtrType, outObjType := getNormalizedObjectTypes(objType)

		require.NotNil(s.T(), outPtrType)
		require.NotNil(s.T(), outObjType)

		assert.Equal(s.T(), reflect.Ptr, outPtrType.Kind())
		assert.NotEqual(s.T(), reflect.Ptr, outObjType.Kind())
	}

	{ // verify pointer input type resolves (non-struct)
		simpleObj := 7
		simpleObjPtr := &simpleObj
		objPtrType := reflect.TypeOf(simpleObjPtr)
		outPtrType, outObjType := getNormalizedObjectTypes(objPtrType)

		require.NotNil(s.T(), outPtrType)
		require.NotNil(s.T(), outObjType)

		assert.Equal(s.T(), reflect.Ptr, outPtrType.Kind())
		assert.NotEqual(s.T(), reflect.Ptr, outObjType.Kind())
	}

	{ // verify multi-depth-pointer input type resolves (non-struct)
		simpleObj := 7
		simpleObjPtr := &simpleObj
		simpleObjPtrPtr := &simpleObjPtr
		simpleObjPtrPtrPtr := &simpleObjPtrPtr

		objPtrType := reflect.TypeOf(simpleObjPtrPtrPtr)
		outPtrType, outObjType := getNormalizedObjectTypes(objPtrType)

		require.NotNil(s.T(), outPtrType)
		require.NotNil(s.T(), outObjType)

		assert.Equal(s.T(), reflect.Ptr, outPtrType.Kind())
		assert.Equal(s.T(), reflect.Int, outObjType.Kind())
	}
}
