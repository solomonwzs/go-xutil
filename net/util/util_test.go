package util

import (
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
