package icmp

import "testing"

func TestIcmp(t *testing.T) {
	raw := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	t.Log(Checksum(raw))
}
