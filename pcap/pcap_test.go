package pcap

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/csignal"
)

func TestPcap(t *testing.T) {
	ch := make(chan os.Signal)
	csignal.Notify(ch, syscall.SIGSEGV)
	go func() {
		for {
			select {
			case sig := <-ch:
				fmt.Println("got signal:", sig)
				os.Exit(1)
			}
		}
	}()

	dev, err := PcapLookupDev()
	if err != nil {
		t.Fatal(err)
	}

	h, err := OpenLive(dev, 512, true, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	if err = h.SetFilter("port 80"); err != nil {
		t.Fatal(err)
	}

	ready, err := h.Wait(5 * time.Second)
	if err != nil {
		t.Fatal(err)
	} else if ready {
		fmt.Println(h.ReadPacket())
	} else {
		fmt.Println("timeout")
	}
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
