package go7941w

type H10301 struct {
	Facility int
	ID       int16
}

// http://www.proxmark.org/files/Documents/125%20kHz%20-%20HID/HID_format_example.pdf
// is the example presented in the proxmark3 source code for properly bit-encoding H10301
// for the T55x7 chips.

func (h H10301) Encode() []byte {
	out := make([]byte, 5)
	// Format:
	// 0x KEY FAC1 FAC2 y0

	out[0] = 0x20
	out[1] = 0x04 ^ ((byte(h.Facility) >> 7) & 0x01)
	out[2] = (byte(h.Facility) << 1) ^ (byte((h.ID >> 13) & 0x01))
	out[3] = byte((h.ID >> 7) & 0xFF)
	out[4] = byte((h.ID<<1)&0xFF) ^ 0x00 // parity bit space

	// Parity bits
	pe := !(parity(int(int(h.Facility<<8) + int((h.ID>>8)&0xF0))))
	if !pe {
		out[1] ^= 0x02
	}
	po := parity(int((int((h.ID>>8)&0x0F) << 8) + int(out[3])))
	if !po {
		out[4] ^= 0x01
	}

	return out
}
