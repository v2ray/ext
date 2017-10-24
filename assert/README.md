# V Assertion Library

This is general assertion library for Golang.

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
