package pcap

import (
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/solomonwzs/goxutil/csignal"
	"github.com/solomonwzs/goxutil/logger"
)

func TestPcap(t *testing.T) {
	ch := make(chan os.Signal)
	csignal.Notify(ch, syscall.SIGSEGV, syscall.SIGABRT)
	go func() {
		for {
			select {
			case sig := <-ch:
				logger.DPrintln("got signal:", sig)
				os.Exit(1)
			}
		}
	}()

	dev, err := PcapFindAllDevs()
	logger.DPrintln(dev)
	if err != nil {
		t.Fatal(err)
	}

	interf, err := net.InterfaceByName(dev[0])
	if err != nil {
		t.Fatal(err)
	}
	h, err := OpenLive(interf, 512, true, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	if err = h.SetFilter("ip[2:2] > 512"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	for i := 0; i < 5; i++ {
		p, _ := h.ReadPacket(false)
		logger.DPrintln(p.Ts, p.Len)
	}
	// ready, err := h.waitForRead(3 * time.Second)
	// if err != nil {
	// 	t.Fatal(err)
	// } else if ready {
	// 	logger.DPrintln(h.ReadPacket(false))
	// } else {
	// 	logger.DPrintln("timeout")
	// }
}

func _BenchmarkPcap(b *testing.B) {
	h, err := OpenLive("eno1", 512, true, 1*time.Second)
	if err != nil {
		b.Fatal(err)
	}
	defer h.Close()
	h.SetFilter("tcp")

	for i := 0; i < b.N; i++ {
		h.ReadPacket(false)
	}
}
