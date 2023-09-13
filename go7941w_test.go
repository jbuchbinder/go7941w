package go7941w

import (
	"bytes"
	"testing"
)

func Test_Go7941W_ParseResponse(t *testing.T) {
	g := Go7941W{}

	testData := []map[string]any{
		{
			// Test case: 0x81 simple reply, parsing 4 bytes (0xdeadbeef) and verify xor checksum
			"in":    []byte{0xcd, 0xdc, 0x00, 0x81, 0x04, 0xde, 0xad, 0xbe, 0xef, 0xa7},
			"out":   []byte{0xde, 0xad, 0xbe, 0xef},
			"error": false,
		},
		{
			// Test case: 0x80 message (error) should trigger a Go error
			"in":    []byte{0xcd, 0xdc, 0x00, 0x80, 0x00, 0xa7},
			"out":   []byte{},
			"error": true,
		},
	}

	for k, v := range testData {
		inV := v["in"].([]byte)
		outV := v["out"].([]byte)
		errV := v["error"].(bool)
		out, err := g.parseResponse(inV)

		if !errV && err != nil {
			t.Errorf("[#%d] parseResponse on %#v produced error %s", k, inV, err.Error())
		} else if errV && err == nil {
			t.Errorf("[#%d] parseResponse on %#v did not produce error", k, inV)
		}
		if bytes.Compare(out, outV) != 0 {
			t.Errorf("[#%d] parseResponse on %#v produced unexpected output %#v", k, inV, outV)
		}
	}
}
