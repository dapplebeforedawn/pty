package pty

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	sys_TIOCGPTN   = 0x4004740F
	sys_TIOCSPTLCK = 0x40045431
)

func open() (pty, tty *os.File, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}

func ptsname(f *os.File) (string, error) {
	var n int
	err := ioctl(f.Fd(), sys_TIOCGPTN, &n)
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(n), nil
}

func ioctl(fd uintptr, cmd uintptr, data *int) error {
	_, _, e := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		cmd,
		uintptr(unsafe.Pointer(data)),
	)
	if e != 0 {
		return syscall.ENOTTY
	}
	return nil
}

func setsize(f *os.File, rows uint16, cols uint16) error {
	return ErrUnsupported
}
