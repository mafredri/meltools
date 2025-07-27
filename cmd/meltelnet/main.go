package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

func loadHistory(l *liner.State, path string) {
	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to open history file %s: %v", path, err)
		}
		return
	}
	defer f.Close()
	l.ReadHistory(f)
}

func saveHistory(l *liner.State, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create history file %s: %v", path, err)
		return
	}
	defer f.Close()
	l.WriteHistory(f)
}

func main() {
	logPath := flag.String("log", "", "path to log file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <host>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	host := flag.Arg(0)
	address := net.JoinHostPort(host, "23")

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(os.Stderr, "Connected to %s\n", address)

	var output io.Writer = os.Stdout
	if *logPath != "" {
		f, err := os.Create(*logPath)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
		defer f.Close()
		output = io.MultiWriter(output, f)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "\nReceived SIGINT...")
		cancel()
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	// Copy server output to stdout and log file.
	copyDone := make(chan struct{})
	go func() {
		_, _ = io.Copy(output, conn)
		close(copyDone)
	}()

	// Replace inputCh and scanner goroutine with liner
	const historyFile = ".meltelnet_history"
	home, _ := os.UserHomeDir()
	histPath := filepath.Join(home, historyFile)
	l := liner.NewLiner()
	defer l.Close()
	l.SetCtrlCAborts(true)
	loadHistory(l, histPath)
	defer saveHistory(l, histPath)

	go func() {
		<-ctx.Done()
		l.Close()
	}()

	for {
		line, err := l.Prompt("")
		if err == liner.ErrPromptAborted || err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Input error: %v", err)
			break
		}
		if strings.TrimSpace(strings.ToLower(line)) == "exit" {
			break
		}
		if line != "" {
			l.AppendHistory(line)
		}
		if _, err := conn.Write([]byte(line + "\r")); err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("Write error: %v", err)
			}
			break
		}
	}

	err = conn.Close()
	if err != nil {
		log.Printf("Error closing connection: %v", err)
	}
	fmt.Fprintln(os.Stderr, "Connection closed")

	<-copyDone
}
