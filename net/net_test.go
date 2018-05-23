package net

import (
	"net"
	"syscall"
	"testing"

	"github.com/solomonwzs/goxutil/net/datalink"
	"github.com/solomonwzs/goxutil/net/ethernet"
	"github.com/solomonwzs/goxutil/net/util"
)

func TestICMP(t *testing.T) {
	dev := "eno1"

	gateway, err := util.GetGateway(dev)
	if err != nil {
		t.Fatal(err)
	}

	hardwareAddr, err := util.GetHardwareAddr(dev, gateway)
	if err == util.ERR_NOT_FOUND {
		hardwareAddr, err = datalink.GetHardwareAddr(dev, gateway, 0)
	}
	if err != nil {
		t.Fatal(err)
	}

	interf, err := net.InterfaceByName(dev)
	if err != nil {
		t.Fatal(err)
	}
	addrs, err := interf.Addrs()
	if err != nil {
		return
	}
	var localIP net.IP = nil
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err != nil {
			t.Fatal(err)
		} else if ipv4 := ip.To4(); ipv4 != nil {
			localIP = ipv4
			break
		}
	}
	if localIP == nil {
		t.Fatal("can not get local ip")
	}

	ethH := &ethernet.EthernetHeader{
		Src:  interf.HardwareAddr,
		Dst:  hardwareAddr,
		Type: syscall.ETH_P_IP,
	}
}
