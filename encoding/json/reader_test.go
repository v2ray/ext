package json_test

import (
	"bytes"
	"testing"

	. "v2ray.com/ext/assert"
	. "v2ray.com/ext/encoding/json"
)

func TestReader(t *testing.T) {
	assert := With(t)

	data := []struct {
		input  string
		output string
	}{
		{
			`
content #comment 1
#comment 2
content 2`,
			`
content content 2`},
		{`content`, `content`},
		{" ", " "},
		{`con/*abcd*/tent`, "content"},
		{`
text // adlkhdf /*
//comment adfkj
text 2*/`, `
text text 2*`},
		{`"//"content`, `"//"content`},
		{`abcd'//'abcd`, `abcd'//'abcd`},
		{`"\""`, `"\""`},
		{`\"/*abcd*/\"`, `\"\"`},
	}

	for _, testCase := range data {
		reader := &Reader{
			Reader: bytes.NewReader([]byte(testCase.input)),
		}

		actual := make([]byte, 1024)
		n, err := reader.Read(actual)
		assert(err, IsNil)
		assert(string(actual[:n]), Equals, testCase.output)
	}
}
