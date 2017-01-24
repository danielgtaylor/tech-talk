package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/googollee/go-socket.io"
	"golang.org/x/crypto/ssh"
)

// Provides a built-in internal SSH mechanism, which works on all operating
// systems. It does _not_ take into account your SSH config (e.g. ~/.ssh/config)
// and so on systems that provide an external SSH client it's better to shell
// out and use your standard config.
func internalSSH(so socketio.Socket) {
	var user string
	var hostPort string

	parts := strings.Split(*sshHost, "@")

	if len(parts) == 1 {
		user = currentUser.Username
		hostPort = parts[0]
	} else {
		user = parts[0]
		hostPort = parts[1]
	}

	if !strings.Contains(hostPort, ":") {
		hostPort = hostPort + ":22"
	}

	var auth []ssh.AuthMethod

	if *key != "" {
		key, err := ioutil.ReadFile(*key)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)

		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}

		auth = append(auth, ssh.PublicKeys(signer))
	}

	if *pass != "" {
		auth = append(auth, ssh.Password(*pass))
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: auth,
	}

	log.Println("Connecting to SSH server...")
	connection, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		log.Printf("Failed to dial: %s", err)
		return
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Printf("Failed to create session: %s", err)
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
		session.Close()
		log.Printf("Request for pseudo terminal failed: %s", err)
		return
	}

	so.On("resize", func(msg map[string]int) {
		log.Printf("Resize: %d cols x %d row\n", msg["col"], msg["row"])

		// The Go SSH implementation doesn't provide this call but does let us
		// send custom commands. See here for a description of the command and
		// structure: https://www.ietf.org/rfc/rfc4254.txt
		data := struct {
			Col uint32
			Row uint32
			W   uint32
			H   uint32
		}{uint32(msg["col"]), uint32(msg["row"]), 0, 0}

		if _, err := session.SendRequest("window-change", false, ssh.Marshal(&data)); err != nil {
			log.Printf("Request for window resize failed: %s", err)
			return
		}
	})

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Printf("Unable to setup stdin for session: %v", err)
		return
	}

	so.On("input", func(msg string) {
		stdin.Write([]byte(msg))
	})

	so.On("disconnection", func() {
		log.Println("Terminal disconnect")
		session.Close()
		connection.Close()
	})

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("Unable to setup stdout for session: %v", err)
		return
	}

	go copyToSocket(stdout, so)

	log.Println("Starting remote shell...")
	err = session.Shell()
	if err != nil {
		log.Printf("Unable to start shell: %v", err)
	}
}
