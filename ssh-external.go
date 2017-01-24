// +build !windows

package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"github.com/googollee/go-socket.io"
	"github.com/kr/pty"
)

const CAN_USE_EXTERNAL = true

// Resize a PTY using system calls. This functionality / utility is missing
// from the kr/pty project so it is added here.
func resizePty(t *os.File, row, col int) error {
	ws := struct {
		ws_row    uint16
		ws_col    uint16
		ws_xpixel uint16
		ws_ypixel uint16
	}{
		uint16(row),
		uint16(col),
		0,
		0,
	}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		t.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

// Connect to an external SSH client (or other external login shell)
func externalSSH(so socketio.Socket) {
	var c *exec.Cmd

	var args []string

	parts := strings.Split(*sshHost, ":")
	if len(parts) == 1 {
		args = append(args, parts[0])
	} else if len(parts) == 2 {
		args = append(args, parts[0])
		args = append(args, "-p")
		args = append(args, parts[1])
	} else {
		log.Panicf("Not sure what to do with host: ", *sshHost)
	}

	// If SSH was explicitly set, then prefer it.
	if *sshHost != DEFAULT_HOST {
		log.Printf("Using external SSH client for %s\n", *sshHost)
		c = exec.Command("/usr/bin/ssh", args...)
	} else {
		// On Mac we should have `/usr/bin/login` which does not require root,
		// so there we use it. Otherwise just start an SSH session with the
		// current user on localhost to get a shell without root.
		if check_access("/usr/bin/login") {
			log.Printf("Using /usr/bin/login for shell as %s\n", currentUser.Username)
			c = exec.Command("/usr/bin/login", "-f", currentUser.Username)
		} else {
			log.Printf("Using external SSH client for %s\n", *sshHost)
			c = exec.Command("/usr/bin/ssh", args...)
		}
	}

	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	so.On("input", func(msg string) {
		f.Write([]byte(msg))
	})

	so.On("resize", func(msg map[string]int) {
		rows, cols, err := pty.Getsize(f)

		if err != nil {
			log.Printf("Error: could not get PTY size. %s\n", err)
			return
		}

		if rows != msg["row"] || cols != msg["col"] {
			log.Printf("Resize: %d cols x %d row\n", msg["col"], msg["row"])
			resizePty(f, msg["row"], msg["col"])
		}
	})

	so.On("disconnection", func() {
		log.Println("Terminal disconnect")
		c.Process.Kill()
	})

	go copyToSocket(f, so)
}
