# MonkeyMock - Powerful. Readable. Bananas.

Easy to use yet powerful mocks library for your Go (golang) tests.
Heavily inspired by RSpec.

Features:
- Human readable DSL that makes mocking both powerful and easy to understand and maintain
- Compatible with stretchr/testify and onsi/ginkgo
- Mock doubles (stand-ins for real objects)
- Partial doubles (layer mocks overtop real objects, only mocking for specific calls)
- Indirect mocks for functions and object methods via monkey patching at runtime ([reference `monkeymock/unsafe` extention](unsafe/README.md))

## Usage

To install MonkeyMock, use go get:
```bash
go get github.com/eshork/monkeymock
```

... or with Go 1.12+ modules, add the require to your `go.mod` file:

```go
require (
  github.com/eshork/monkeymock
)
```

These will make the following packages available to your Go imports:
```
github.com/eshork/monkeymock
github.com/eshork/monkeymock/unsafe
```

Import the eshork/monkeymock package into your code using this template:
```go
package yours

import (
  "testing"
  "github.com/eshork/monkeymock"
)

func TestSomething(t *testing.T) {
  // set up a mock atop an existing object
  mockObj := monkeymock.Expect(yourObj).
    ToReceive("method"). // object method name (string symbol value)
    Once(). // how many times will this be called? once (obviously)
    WithArgs("string_arg"). // expect specific arg values
    AndCallOriginal(). // preserve the original functionality
    WithReturns(7) // expect a specific return value

  // do some test actions using the mock
  mockObject.method("somevalue")

  // assert that all established expectations have been met
  monkeymock.AssertExpectations(t)

  // this specific example will fail with an args-based assertion failure!
}
```

For more examples, see: EXAMPLES.md


## This not possible without
- [RSpec](https://relishapp.com/rspec) - Thanks for the great inspiration
- [bouk/monkey](https://github.com/bouk/monkey) - Thanks for making a great runtime monkey patch library for Go
