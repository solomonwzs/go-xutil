package csignal

/*
#include <signal.h>
#include <stdio.h>
#include <sys/socket.h>
#include <unistd.h>

#define SIZEOF_INT sizeof(int)

int csock;

void
_notify_go(int sig, siginfo_t *info, void *context) {
	printf("c: it is evil!\n");
	write(csock, &sig, sizeof(sig));
}

int
_sigaltstack() {
	static char _stack[SIGSTKSZ];
	stack_t ss = {
		.ss_size = SIGSTKSZ,
		.ss_sp = _stack,
	};

	if (sigaltstack(&ss, 0) != 0) {
		perror("sigaltstack: ");
		return -1;
	}
	return 0;
}

int
_signal_handle(int flags, int sig) {
	struct sigaction action;
	sigemptyset(&action.sa_mask);
	action.sa_sigaction = _notify_go;
	action.sa_flags = flags;

	if (sigaction(sig, &action, NULL) != 0) {
		perror("csignal: ");
		return -1;
	}
	return 0;
}
*/
import "C"
import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"
	"syscall"
)

const (
	SA_NOCLDSTOP = C.SA_NOCLDSTOP
	SA_ONSTACK   = C.SA_ONSTACK
	SA_RESETHAND = C.SA_RESETHAND
	SA_RESTART   = C.SA_RESTART
	SA_SIGINFO   = C.SA_SIGINFO
	SA_NOCLDWAIT = C.SA_NOCLDWAIT
	SA_NODEFER   = C.SA_NODEFER
)

var handlers struct {
	sync.Mutex
	m map[os.Signal][]chan<- os.Signal
}

func init() {
	fd, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}
	C.csock = C.int(fd[0])
	// C._sigaltstack()

	handlers.m = map[os.Signal][]chan<- os.Signal{}

	goSock := os.NewFile(uintptr(fd[1]), "goxutil-csignal")
	go recv_loop(goSock)
}

func recv_loop(goSock *os.File) {
	p := make([]byte, C.SIZEOF_INT)
	for {
		if _, err := goSock.Read(p); err == nil {
			sig := syscall.Signal(binary.LittleEndian.Uint32(p))
			handlers.Lock()
			if chs, exist := handlers.m[sig]; exist {
				for _, ch := range chs {
					select {
					case ch <- sig:
					default:
					}
				}
			}
			handlers.Unlock()
		} else if err == io.EOF || err == os.ErrClosed {
			return
		}
	}
}

func catchSignal(flags int, sig syscall.Signal) (err error) {
	if C._signal_handle(C.int(flags), C.int(sig)) != 0 {
		err = errors.New("[csignal] set signal action fail")
	}
	return
}

func Notify(c chan<- os.Signal, sig ...syscall.Signal) {
	handlers.Lock()
	defer handlers.Unlock()

	for _, s := range sig {
		if _, exist := handlers.m[s]; !exist {
			if err := catchSignal(
				SA_NOCLDSTOP|SA_SIGINFO|SA_ONSTACK, s); err == nil {
				handlers.m[s] = []chan<- os.Signal{c}
			}
		} else {
			handlers.m[s] = append(handlers.m[s], c)
		}
	}
}
