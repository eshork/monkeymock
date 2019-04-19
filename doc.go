/*
Package monkeymock provides a flexible way to mock your objects and verify
calls are happening as expected.

## General usage

Set up one more more Mock expectations

Execute your example scenario

Evaluate the overall


## Scopes

All monkeymock constructs are intended to have a per-example lifecycle.
In the typical use case, the final call to `monkeymock.AssertExpectations(t)`
will clear all expectations that been configured thus far.

If this workflow does not fit your use case, see [AssertExpectations] for
details on how to preserve state across multiple calls for general assertion.



Message expectations are verified
after each example. Doubles, method stubs, stubbed constants, etc. are all cleaned up after
each example. This ensures that each example can be run in isolation, and in any order.

It is perfectly fine to set up doubles, stubs, and message expectations in a
before(:example) hook, as that hook is executed in the scope of the example:
*/
package monkeymock
