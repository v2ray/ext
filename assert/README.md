# V Assertion Library

Yet another assertion library for Golang.

## Quick Start

```go
import (
    . "v2ray.com/ext/assert"
)

func TestStringEquals(t *testing.T) {
    assert := With(t)

    str := "Hello" + " " + "World"
    assert(str, Equals, "Hello World")
}
```

## Usage

The `assert` function take at least 2 parameters: `assert(value, matcher, expectations...)`, where:
* `value` is the value to be asserted.
* `matcher` is a "function" to assert the value.
* `expectations` are optional parameters that will be passed into `matcher` for the assertion.

There are several predefined matchers:
* Equals
* LessThan
* GreaterThan
* AtLeast
* AtMost
* IsNil
* IsNotNil
* HasSubstring
