// +build !windows

package pty

import (
	"os"
	"syscall"
	"unsafe"
)

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func Getsize(t *os.File) (rows, cols int, err error) {
	var ws winsize
	err = windowrect(&ws, t.Fd(), syscall.TIOCGWINSZ)
	return int(ws.ws_row), int(ws.ws_col), err
}

func Setsize(t *os.File, rows, cols uint16) error {
	ws := winsize{
		ws_row: rows,
		ws_col: cols,
		ws_xpixel: 0,
		ws_ypixel: 0,
	}
  return windowrect(&ws, t.Fd(), syscall.TIOCSWINSZ)
}

type winsize struct {
	ws_row    uint16
	ws_col    uint16
	ws_xpixel uint16
	ws_ypixel uint16
}

func windowrect(ws *winsize, fd, call uintptr) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		call,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}