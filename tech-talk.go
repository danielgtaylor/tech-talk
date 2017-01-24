package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/googollee/go-socket.io"
)

type TemplateValues struct {
	Prefix   string
	Markdown string
}

const DEFAULT_HOST = "localhost"
const VERSION = "1.2.0"

var indexTemplate *template.Template
var socketServer *socketio.Server
var mdFilename string
var currentUser *user.User

var sshType *string
var sshHost *string
var key *string
var pass *string
var noBrowser *bool

// Checks if a file exists and can be accessed.
func check_access(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Copy data from a reader (e.g. PTY or Stdin pipe) to the web socket.
func copyToSocket(r io.Reader, so socketio.Socket) {
	for {
		data := make([]byte, 512)
		n, err := r.Read(data)
		if err != nil {
			log.Println(err)
			break
		}
		if n > 0 {
			so.Emit("output", string(data))
		}
	}
}

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

		if *sshType == "internal" || CAN_USE_EXTERNAL == false {
			log.Println("Using internal SSH client")
			internalSSH(so)
		} else {
			externalSSH(so)
		}
	})

	socketServer.On("error", func(so socketio.Socket, err error) {
		log.Println("Error:", err)
	})
}

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal("Couldn't get current user!")
	}
	currentUser = u

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tech-talk [slides.md]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	// Connection options
	sshHost = flag.String("host", DEFAULT_HOST, "SSH connection [user@]`hostname`[:port]")
	sshType = flag.String("ssh", "auto", "SSH `method` [auto, internal]")

	// Auth options
	keyDefault := ""

	idRsaPath := path.Join(currentUser.HomeDir, ".ssh", "id_rsa")
	if check_access(idRsaPath) {
		keyDefault = idRsaPath
	}

	key = flag.String("key", keyDefault, "SSH private key `path` (for internal SSH)")
	pass = flag.String("pass", "", "SSH `password` (for internal SSH)")

	// Misc options
	noBrowser = flag.Bool("n", false, "Do not automatically open browser")
	version := flag.Bool("v", false, "Alias for --version")
	flag.BoolVar(version, "version", false, "Print program version and exit")

	flag.Parse()
	args := flag.Args()

	if *version {
		fmt.Printf("tech-talk: %s\n", VERSION)
		return
	}

	if len(args) > 0 {
		if !check_access(args[0]) {
			log.Fatalf("Cannot access %s", args[0])
		}

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

	log.Println("Server started on http://localhost:4000/")

	if !*noBrowser && check_access("/usr/bin/open") {
		c := exec.Command("/usr/bin/open", "http://localhost:4000")
		c.Start()
	}

	log.Panic(s.ListenAndServe())
}
