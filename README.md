# Tech Talk

An opinionated Markdown-based technical slideshow tool with a built-in terminal that just works.

![techtalk](https://cloud.githubusercontent.com/assets/106826/22236695/229d6c66-e1bc-11e6-9bf2-5b8cfa500e91.gif)

- Write your slides with [Markdown](https://github.com/gnab/remark/wiki/Markdown)
- Built-in terminal (both local or via SSH)
- Supports presenter mode with notes & timer, cloned displays, etc.
- Simple, self-contained executable for Winows, Mac & Linux
- Trivial to customize and distribute within your company

## Usage

Start by downloading a [release](https://github.com/danielgtaylor/tech-talk/releases) or your company's variant. You probably want to put it into `/usr/local/bin/tech-talk` on Mac / Linux. On Windows you can save `tech-talk.exe` anywhere as it is fully self-contained / portable.

```sh
# See a demo slideshow
tech-talk

# Use your own slides
tech-talk slides.md

# Use SSH to connect to some server for the terminal
tech-talk -host user@hostname slides.md

# Custom ports are also allowed
tech-talk -host user@hostname:port slides.md

# Windows users must pass an authentication method as an internal SSH
# mechanism is used instead of OpenSSH. Keys are recommended over passwords.
tech-talk -host user@hostname -pass cleartext-password slides.md
tech-talk -host user@hostname -key id_rsa slides.md
tech-talk -host user@hostname -key my-key.pem slides.md

# Mac / Linux users can force the use of the internal SSH client and use the
# same options that Windows users would use from above.
tech-talk -ssh internal ...
```

Then go to [http://localhost:4000/](http://localhost:4000/) to view.

## Customization

Sending around boilerplate or configs to everyone sucks. Build your own self-contained executable with your company's theme and let people focus on making great talks.

First you'll want to install [Go](https://golang.org/) and [Glide](https://github.com/Masterminds/glide) (e.g. `curl https://glide.sh/get | sh`). Then:

```sh
mkdir -p $GOPATH/src/github.com/danielgtaylor
cd $GOPATH/src/github.com/danielgtaylor
git clone https://github.com/danielgtaylor/tech-talk.git
cd tech-talk
$GOPATH/bin/glide install
```

Now, make your customizations!

- [Slideshow prefix (prepended to first slide)](https://github.com/danielgtaylor/tech-talk/tree/master/data/prefix.md)
- [Stylesheet (fonts, colors, layout, etc)](https://github.com/danielgtaylor/tech-talk/tree/master/www/style.css)
- [Javascript (transitions)](https://github.com/danielgtaylor/tech-talk/tree/master/www/script.js)
- [Terminal default font size](https://github.com/danielgtaylor/tech-talk/tree/master/www/wetty/wetty.js)
- [Example slideshow](https://github.com/danielgtaylor/tech-talk/tree/master/data/example.md)
- [HTML template](https://github.com/danielgtaylor/tech-talk/blob/master/data/index.template)

Once you are ready:

```sh
# Build for your OS and test
./build.sh
./tech-talk

# Or, cross-compile, e.g.
GOOS=linux GOARCH=386 ./build.sh
```

Remember to run `./build.sh` each time you make a change, and your browser may cache items so Cmd+Shift+R or Ctrl+Shift+R to force a refresh are useful.

Now you can upload your executables somewhere like Google Drive and share them within the company.

## Acknowledgments

This project is possible because of the amazing work done by many people in the following projects, many of which are used with slight modifications or custom settings:

- [Remark.js](https://github.com/gnab/remark#remark)
- [Socket.io](https://github.com/googollee/go-socket.io)
- [pty](https://github.com/kr/pty)
- [Wetty](https://github.com/krishnasrinivas/wetty)
- [go-bindata](https://github.com/jteeuwen/go-bindata)
- [crypto/ssh](https://godoc.org/golang.org/x/crypto/ssh)

So, how does this beast work? Simple, really. It starts both a web server and Socket.IO server, renders an index page with Remark.js using the user-supplied Markdown, which in turn contains an `iframe` with a terminal that uses the socket to communicate with a PTY running either a login shell or `ssh` session. For the internal SSH client there is no client-side PTY but instead a couple pipes through an SSH connection.

## License

https://dgt.mit-license.org/
