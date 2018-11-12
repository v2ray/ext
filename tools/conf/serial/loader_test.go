package serial_test

import (
	"bytes"
	"strings"
	"testing"

	"v2ray.com/ext/tools/conf/serial"
)

func TestLoaderError(t *testing.T) {
	testCases := []struct {
		Input  string
		Output string
	}{
		{
			Input: `{
				"log": {
					// abcd
					0,
					"loglevel": "info"
				}
		}`,
			Output: "line 4 char 7",
		},
		{
			Input: `{
				"log": {
					// abcd
					"loglevel": "info",
				}
		}`,
			Output: "line 6 char 0",
		},
		{
			Input: `{
				"port": 1,
				"inbounds": [{
					"protocol": "test"
				}]
		}`,
			Output: "parse json config",
		},
	}
	for _, testCase := range testCases {
		reader := bytes.NewReader([]byte(testCase.Input))
		_, err := serial.LoadJSONConfig(reader)
		errString := err.Error()
		if !strings.Contains(errString, testCase.Output) {
			t.Error("unexpected output from json: ", testCase.Input, ". expected ", testCase.Output, ", but actualy ", errString)
		}
	}
}
