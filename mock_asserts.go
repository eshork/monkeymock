package monkeymock

import (
	"testing"
)

type mockAssertions interface {
	AssertExpectations(t *testing.T, opts ...interface{})
}

// AssertExpectations for a single Mock.
func (m *mockStruct) AssertExpectations(t *testing.T, opts ...interface{}) {
	t.Helper()
	m.assertMethods(t)
}

/// General module-level assertions

// AssertExpectations across all Mock instances.
func AssertExpectations(t *testing.T, opts ...interface{}) {
	//
}

// ClearExpectations resets the board, removing all existing expectations for every Mock.
func ClearExpectations(t *testing.T, opts ...interface{}) {
	//
}
