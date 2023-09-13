package go7941w

func parity(n int) bool {
	p := false
	for n != 0 {
		if n&1 != 0 {
			p = !p
		}
		n = n >> 1
	}
	return p
}
