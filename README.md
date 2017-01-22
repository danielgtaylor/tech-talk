# Tech Talk

An opinionated Markdown-based slideshow tool that just works.

![techtalk](https://cloud.githubusercontent.com/assets/106826/22179439/851145a8-e008-11e6-9c31-c391c4025546.gif)

- Write your slides with Markdown
- Built-in terminal (both local or via SSH)
- Simple, self-contained executable for Windows, Mac, Linux
- Trivial to customize and distribute within your company

## Usage

Start by downloading a [release](https://github.com/danielgtaylor/tech-talk/releases) or your company's variant. You probably want to put it into `/usr/local/bin/tech-talk` on Mac / Linux.

```sh
# See a demo slideshow
tech-talk

# Use your own slides
tech-talk my-cool-slides.md

# Use SSH to connect to some server for the terminal
tech-talk -host user@hostname my-cool-slides.md
```

Then go to [http://localhost:4000/](http://localhost:4000/) to view.

## Customization

Sending around boilerplate or configs to everyone sucks. Build your own self-contained executable with your company's theme and let people focus on making great talks.

First you'll want to install [Go](). Then:

```sh
git clone https://github.com/danielgtaylor/tech-talk.git
cd tech-talk
export GOPATH=`pwd`
go get -u github.com/jteeuwen/go-bindata/...
```

Now, make your customizations!

- [Slideshow prefix (prepended to first slide)](https://github.com/danielgtaylor/tech-talk/tree/master/data/prefix.md)
- [Stylesheet (fonts, colors, layout, etc)](https://github.com/danielgtaylor/tech-talk/tree/master/www/style.css)
- [Javascript (transitions)](https://github.com/danielgtaylor/tech-talk/tree/master/www/script.js)
- [Terminal default font size](https://github.com/danielgtaylor/tech-talk/tree/master/www/wetty/wetty.js)
- [Example slideshow](https://github.com/danielgtaylor/tech-talk/tree/master/www/style.css)

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

So, how does this beast work? Simple, really. It starts both a web server and Socket.IO server, renders an index page with Remark.js using the user-supplied Markdown, which in turn contains an `iframe` with a terminal that uses the socket to communicate with a PTY running either a login shell or `ssh` session.

## License

https://dgt.mit-license.org/
