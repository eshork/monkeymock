package testsupports

// Tests that we can only run from an internal perspective

import (
	// 	"errors"
	// 	"fmt"
	// 	"regexp"
	// 	"runtime"
	// 	"sync"

	"fmt"
	"reflect"

	// 	"time"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// AssertNotSame ...
func AssertNotSame(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	expectedPtr, actualPtr := reflect.ValueOf(expected), reflect.ValueOf(actual)
	if expectedPtr.Kind() != reflect.Ptr || actualPtr.Kind() != reflect.Ptr {
		return assert.Fail(t, "Invalid operation: both arguments must be pointers", msgAndArgs...)
	}

	expectedType, actualType := reflect.TypeOf(expected), reflect.TypeOf(actual)
	if expectedType != actualType {
		return assert.Fail(t, fmt.Sprintf("Pointer expected to be of type %v, but was %v",
			expectedType, actualType), msgAndArgs...)
	}

	if expected == actual {
		return assert.Fail(t, fmt.Sprintf("Same: \n"+
			"expected  : %p %#v\n"+
			"to not be : %p %#v", expected, expected, actual, actual), msgAndArgs...)
	}

	return true
}

// AssertSame ...
func AssertSame(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	expectedPtr, actualPtr := reflect.ValueOf(expected), reflect.ValueOf(actual)
	if expectedPtr.Kind() != reflect.Ptr || actualPtr.Kind() != reflect.Ptr {
		return assert.Fail(t, "Invalid operation: both arguments must be pointers", msgAndArgs...)
	}

	expectedType, actualType := reflect.TypeOf(expected), reflect.TypeOf(actual)
	if expectedType != actualType {
		return assert.Fail(t, fmt.Sprintf("Pointer expected to be of type %v, but was %v",
			expectedType, actualType), msgAndArgs...)
	}

	if expected != actual {
		return assert.Fail(t, fmt.Sprintf("Not same: \n"+
			"expected: %p %#v\n"+
			"actual  : %p %#v", expected, expected, actual, actual), msgAndArgs...)
	}

	return true
}

type tHelper interface {
	Helper()
}
