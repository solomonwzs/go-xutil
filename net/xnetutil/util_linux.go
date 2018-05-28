package xnetutil

/*
#include <linux/route.h>
*/
import "C"
import (
	"bufio"
	"errors"
	"net"
	"os"
	"regexp"
	"strconv"
)

const (
	_ROUTE_TABLE = "/proc/net/route"
	_ARP_TABLE   = "/proc/net/arp"
)

var (
	ERR_NOT_FOUND = errors.New("[net] not found")
)

func GetGateway(dev string) (addr net.IP, err error) {
	file, err := os.Open(_ROUTE_TABLE)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile("[^ \t]+")
	for scanner.Scan() {
		line := scanner.Text()
		arr := re.FindAllString(line, -1)
		if len(arr) != 11 || arr[0] != dev {
			continue
		}

		flags, err0 := strconv.Atoi(arr[3])
		if err0 != nil {
			return nil, err0
		}
		if flags&C.RTF_UP == 0 || flags&C.RTF_GATEWAY == 0 {
			continue
		}

		gateway := arr[2]
		size := len(gateway) / 2
		addr = make([]byte, size, size)
		for i := 0; i < size; i++ {
			j, err1 := strconv.ParseUint(gateway[i*2:i*2+2], 16, 8)
			if err1 != nil {
				return nil, err1
			}
			addr[size-i-1] = byte(j)
		}
		return addr, err
	}
	return nil, ERR_NOT_FOUND
}

func GetHardwareAddr(dev string, ip net.IP) (
	hw net.HardwareAddr, err error) {
	file, err := os.Open(_ARP_TABLE)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ipStr := ip.String()
	re := regexp.MustCompile("[^ \t]+")
	for scanner.Scan() {
		line := scanner.Text()
		arr := re.FindAllString(line, -1)
		if len(arr) != 6 || arr[5] != dev || arr[0] != ipStr {
			continue
		}

		return net.ParseMAC(arr[3])
	}
	return nil, ERR_NOT_FOUND
}
