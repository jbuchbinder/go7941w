package go7941w

import (
	"encoding/hex"
	"testing"
)

func Test_H10301_Encode(t *testing.T) {
	h := H10301{
		Facility: 0x63,
		ID:       0x1C25,
	}
	o := h.Encode()
	t.Logf("h (encoded) : %#v", o)
	//n := (o[0] << 32) + (o[1] << 24) + (o[2] << 16) + (o[3] << 8) + o[4]
	if hex.EncodeToString(o) != "2006c6384a" {
		t.Errorf("hex(): %s != '2006c6384a", hex.EncodeToString(o))
	}
}
