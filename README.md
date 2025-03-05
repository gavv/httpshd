# httpshd

[![Build](https://github.com/gavv/httpshd/workflows/build/badge.svg)](https://github.com/gavv/httpshd/actions)

`httpshd` is a tiny HTTP server that reads shell command from request body, executes it, and sends its output via HTTP response.

There are two features for better support of interactivity:

* command is executed inside a PTY
* command output is not buffered, chunked transfer encoding is used

Thanks to this, the output of a long-running command is delivered to the client immediately.

## Security warning

This program is absolutely insecure!

There is neither encryption, nor authentication. Anyone can execute any command with the permissions of the user which is running `httpshd`. Hence, it's recommended to use it only within private networks or tunnels.

As suggested on reddit, you can use SSH port forwarding to secure it (see example below).

## Installation

First, install Go.

On Debian / Ununtu:

```
sudo apt install golang-go
```

On macOS:

```
brew install go
```

Then run this command, which will download, build, and install `httpshd` executable into `$GOPATH/bin` (it's `~/go/bin` if `GOPATH` environment variable is not set):

```
go install -v github.com/gavv/httpshd@latest
```

Alternatively, you can install from sources:

```
git clone https://github.com/gavv/httpshd.git
cd httpshd
go build
./httpshd -help
```

## Usage example

Server:

```
$ ./httpshd
21:00:03.585 starting server at 0.0.0.0:3333
21:00:14.734 [127.0.0.1:37786] executing command: /bin/zsh -c "ls -l"
21:00:14.741 [127.0.0.1:37786] success
...
```

Client:

```
$ echo "ls -l" | curl -XPOST -d @- -s http://<server>:3333
total 6256
-rw-r--r-- 1 user user     149 Jun 16 21:01 go.mod
-rw-r--r-- 1 user user     370 Mar 24 19:21 go.sum
-rw-r--r-- 1 user user    1080 Jun 16 21:15 LICENSE
-rw-r--r-- 1 user user    1673 Jun 16 20:59 main.go
-rw-r--r-- 1 user user    1235 Jun 16 21:14 README.md
```

## Using with SSH port forwarding

Server:

```
$ ./httpshd -host 127.0.0.1
21:00:03.585 starting server at 127.0.0.1:3333
21:00:14.734 [127.0.0.1:37786] executing command: /bin/zsh -c "ls -l"
21:00:14.741 [127.0.0.1:37786] success
...
```

Client:

```
$ ssh -fN -L 127.0.0.1:3333:127.0.0.1:3333 <server>

$ echo "ls -l" | curl -XPOST -d @- -s http://127.0.0.1:3333
total 6256
-rw-r--r-- 1 user user     149 Jun 16 21:01 go.mod
-rw-r--r-- 1 user user     370 Mar 24 19:21 go.sum
-rw-r--r-- 1 user user    1080 Jun 16 21:15 LICENSE
-rw-r--r-- 1 user user    1673 Jun 16 20:59 main.go
-rw-r--r-- 1 user user    1235 Jun 16 21:14 README.md
```

## Options

```
$ ./httpshd -help
Usage of ./httpshd
  -help
    	print help
  -host string
    	interface address (default "0.0.0.0")
  -port int
    	port number (default 3333)
  -sh string
    	path to shell (default "/bin/zsh")
```

## Motivation

I created this program because I needed to run compilation on macOS remotely from my Linux box. Unfortunately, Xcode toolchain has many difficulties with SSH, and some build steps like app signing just didn't work even with some well-known workarounds.

Meanwhile, just sending commands to `nc` connected to `sh` worked, given that `nc` is started inside an interactive session (i.e. you login, open terminal, and run it). The reason is that `nc`, as well as `httpshd`, and unlike `sshd`, doesn't create a new user session, which is handled specially by the OS.

Since `nc` is not very handy for this specific use case, this tool was created.

## Authors

See [here](https://github.com/gavv/httpshd/graphs/contributors).

## License

[MIT](LICENSE)
