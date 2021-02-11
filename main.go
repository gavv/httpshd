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
	"strings"
	"time"

	"github.com/creack/pty"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s %s", time.Now().Format("15:04:05.000"), bytes)
}

type unbufferedWriter struct {
	http.ResponseWriter
}

func (w unbufferedWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)

	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}

	return n, err
}

func main() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})

	defaultShell := os.Getenv("SHELL")
	if defaultShell == "" {
		defaultShell = "/bin/sh"
	}

	host := flag.String("host", "0.0.0.0", "interface address")
	port := flag.Int("port", 3333, "port number")
	shell := flag.String("sh", defaultShell, "path to shell")
	help := flag.Bool("help", false, "print help")

	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("%s: [%s] can't read body", req.RemoteAddr, err.Error())
			return
		}

		s := string(b)
		s = strings.TrimSpace(s)

		log.Printf("[%s] executing command: %s -c %q", req.RemoteAddr, *shell, s)

		cmd := exec.Command(*shell, "-c", s)

		fp, err := pty.Start(cmd)
		if err != nil {
			log.Printf("[%s] command failed: %s", req.RemoteAddr, err.Error())
			return
		}

		io.Copy(unbufferedWriter{w}, fp)

		log.Printf("[%s] success", req.RemoteAddr)
	})

	address := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("starting server at %s\n", address)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Printf("can't start server: %s", err.Error())
	}
}
