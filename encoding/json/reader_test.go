package json_test

import (
	"bytes"
	"io"
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

func TestReader1(t *testing.T) {
	assert := With(t)

	type dataStruct struct {
		input  string
		output string
	}

	bufLen := 8

	data := []dataStruct{
		{"loooooooooooooooooooooooooooooooooooooooog", "loooooooooooooooooooooooooooooooooooooooog"},
		{`{"t": "\/testlooooooooooooooooooooooooooooong"}`, `{"t": "\/testlooooooooooooooooooooooooooooong"}`},
		{`{"t": "\/test"}`, `{"t": "\/test"}`},
		{`"\// fake comment"`, `"\// fake comment"`},
		{`"\/\/\/\/\/"`, `"\/\/\/\/\/"`},
	}

	for _, testCase := range data {
		reader := &Reader{
			Reader: bytes.NewReader([]byte(testCase.input)),
		}
		target := make([]byte, 0)
		buf := make([]byte, bufLen)
		var n int
		var err error
		for n, err = reader.Read(buf); err == nil; n, err = reader.Read(buf) {
			assert(n, AtMost, len(buf))
			target = append(target, buf[:n]...)
			buf = make([]byte, bufLen)
		}
		if err == io.EOF {
			assert(string(target), Equals, testCase.output)
		} else {
			assert(err, IsNil)
			assert(string(target), Equals, testCase.output)
		}
	}

}
