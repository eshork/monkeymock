package monkeymock

// Tests that we can only run from an internal perspective

import (
	. "github.com/eshork/monkeymock/testsupports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *testInternalMocks) TestMockListManipilation() {
	// internal mock list size should increment with each Expect statement
	assert.Equal(s.T(), 0, sizeOfMockList())
	var _ = Expect(struct{}{}) // mock of an empty anonymous struct
	assert.Equal(s.T(), 1, sizeOfMockList())
	var _ = Expect(struct{}{}) // mock of an empty anonymous struct
	assert.Equal(s.T(), 2, sizeOfMockList())

	// clearing the list should make it zero-length (obviously)
	clearMockList()
	assert.Equal(s.T(), 0, sizeOfMockList())

	// can remove the first mock
	{
		first := Expect(struct{}{})  // mock of an empty anonymous struct
		second := Expect(struct{}{}) // mock of an empty anonymous struct
		third := Expect(struct{}{})  // mock of an empty anonymous struct
		require.Equal(s.T(), 3, sizeOfMockList())
		removeFromMockList(first)
		require.Equal(s.T(), 2, sizeOfMockList())
		AssertSame(s.T(), gTheMockList[0], second)
		AssertSame(s.T(), gTheMockList[1], third)
		clearMockList()
	}

	// can remove the last mock
	{
		clearMockList()
		first := Expect(struct{}{})  // mock of an empty anonymous struct
		second := Expect(struct{}{}) // mock of an empty anonymous struct
		third := Expect(struct{}{})  // mock of an empty anonymous struct
		require.Equal(s.T(), 3, sizeOfMockList())
		removeFromMockList(third)
		require.Equal(s.T(), 2, sizeOfMockList())
		AssertSame(s.T(), gTheMockList[0], first)
		AssertSame(s.T(), gTheMockList[1], second)
		clearMockList()
	}
	// can remove the middle mock
	{
		clearMockList()
		first := Expect(struct{}{})  // mock of an empty anonymous struct
		second := Expect(struct{}{}) // mock of an empty anonymous struct
		third := Expect(struct{}{})  // mock of an empty anonymous struct
		require.Equal(s.T(), 3, sizeOfMockList())
		removeFromMockList(second)
		require.Equal(s.T(), 2, sizeOfMockList())
		AssertSame(s.T(), gTheMockList[0], first)
		AssertSame(s.T(), gTheMockList[1], third)
		clearMockList()
	}
}
