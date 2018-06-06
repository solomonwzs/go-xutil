package pcap

import (
	"fmt"
	"testing"
	"time"
)

func TestPcap(t *testing.T) {
	fmt.Println(PcapLookupDev())
	h, err := OpenLive("eno1", 512, true, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	fmt.Println(h.SetFilter("tcp"))
	fmt.Println(h.ReadPacket())
}

func _BenchmarkPcap(b *testing.B) {
	h, err := OpenLive("eno1", 512, true, 1*time.Second)
	if err != nil {
		b.Fatal(err)
	}
	defer h.Close()
	h.SetFilter("tcp")

	for i := 0; i < b.N; i++ {
		h.ReadPacket()
	}
}
