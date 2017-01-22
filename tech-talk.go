package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/googollee/go-socket.io"
	"github.com/kr/pty"
)

type TemplateValues struct {
	Prefix   string
	Markdown string
}

var indexTemplate *template.Template
var socketServer *socketio.Server
var mdFilename string

var sshHost *string

// Return an HTML page with the slideshow
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// If we aren't getting the index itself, serve static files in the
	// same directory as the input Markdown slides file.
	if r.URL.Path != "/" {
		http.FileServer(http.Dir(filepath.Dir(mdFilename))).ServeHTTP(w, r)
		return
	}

	var data TemplateValues

	// Read the file on each request so that updates get applied when working
	// on the slideshow.
	var b []byte

	if mdFilename != "" {
		b, _ = ioutil.ReadFile(mdFilename)
		data.Markdown = string(b)
	} else {
		b, _ = Asset("data/example.md")
		data.Markdown = string(b)
	}

	b, _ = Asset("data/prefix.md")
	data.Prefix = string(b)

	w.Header().Add("Content-Type", "text/html")
	indexTemplate.Execute(w, data)
}

// Create a new socket server to handle communication with a PTY shell.
// This allows you to run stuff in a terminal without ever leaving the
// slideshow.
func createSocketServer() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	socketServer = server

	socketServer.On("connection", func(so socketio.Socket) {
		log.Printf("Terminal connected from %s\n", so.Request().RemoteAddr)

		var c *exec.Cmd

		// If SSH was explicitly set, then prefer it.
		if *sshHost != "localhost" {
			c = exec.Command("/usr/bin/ssh", *sshHost)
		} else {
			// On Mac we should have `/usr/bin/login` which does not require root,
			// so there we use it. Otherwise just start an SSH session with the
			// current user on localhost to get a shell without root.
			if _, err := os.Stat("/usr/bin/login"); err == nil {
				c = exec.Command("/usr/bin/login", "-f", os.Getenv("USER"))
			} else {
				c = exec.Command("/usr/bin/ssh", *sshHost)
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
				pty.Setsize(f, uint16(msg["row"]), uint16(msg["col"]))
			}
		})

		so.On("disconnection", func() {
			log.Println("Terminal disconnect")
		})

		go func() {
			// Read from the PTY until we can't read anymore!
			for {
				data := make([]byte, 512)
				n, err := f.Read(data)
				if err != nil {
					log.Println(err)
					break
				}
				if n > 0 {
					so.Emit("output", string(data))
				}
			}
		}()
	})

	socketServer.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tech-talk [slides.md]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	sshHost = flag.String("host", "localhost", "SSH hostname")

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		mdFilename = args[0]
	}

	// Start web sockets
	createSocketServer()
	http.Handle("/wetty/socket.io/", socketServer)

	// Setup web server
	indexBytes, _ := Asset("data/index.template")
	indexTemplate = template.Must(template.New("index").Parse(string(indexBytes)))

	http.HandleFunc("/", indexHandler)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(
				&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "www"})))

	s := &http.Server{
		Addr:           ":4000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// TODO: run `open http://localhost:4000/`
	log.Println("Server started on http://localhost:4000/")

	log.Panic(s.ListenAndServe())
}
