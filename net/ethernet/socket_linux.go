package ethernet

/*
#include <sys/socket.h>
#include <sys/ioctl.h>
#include <sys/time.h>

#include <asm/types.h>

#include <errno.h>
#include <math.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <arpa/inet.h>

#include <linux/if_packet.h>
#include <linux/if_ether.h>
#include <linux/if_arp.h>

int
t_socket() {
	return socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL));
}

int
ioctl_ifr(int s, char *dev) {
	struct ifreq ifr;
	strncpy(ifr.ifr_name, dev, IFNAMSIZ);
	if (ioctl(s, SIOCGIFINDEX, &ifr) == -1) {
		return errno;
	}
	if (ioctl(s, SIOCGIFHWADDR, &ifr) == -1) {
		return errno;
	}
	return 0;
}
*/
import "C"
import "syscall"

func IoctlIfr(fd int, dev string) (err error) {
	errno := int(C.ioctl_ifr(C.int(fd), C.CString(dev)))
	if errno == 0 {
		return nil
	} else {
		return syscall.Errno(errno)
	}
}

func TSocket() (fd int, err error) {
	fd = int(C.t_socket())
	return
}
