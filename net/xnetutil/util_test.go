package xnetutil

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	ip, err := GetGateway("eno1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ip)
	fmt.Println(GetHardwareAddr("eno1", ip))
}

func TestChecksum(t *testing.T) {
	buf := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	fmt.Println(Checksum(buf))

	c := NewChecksumer()
	c.Write(buf[:5])
	fmt.Println(c.Sum(buf[5:]))
	fmt.Println(c.Sum(buf[5:]))

	h := hmac.New(md5.New, nil)
	fmt.Println(h.BlockSize())
	h.Write(buf[:5])
	fmt.Println(h.Sum(buf[5:]))
	fmt.Println(h.Sum(buf[5:]))
	fmt.Println(h.Size())
}
