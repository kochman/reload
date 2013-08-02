package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Give me the path to a Go file to reload upon modification")
		os.Exit(1)
	}

	filepath := flag.Args()[0]
	args := flag.Args()[1:]

	watchForChange(filepath, args)
}

func runGoFile(filepath string, args []string, done chan bool) {
	// convert arg slice to string
	argString := ""
	for i, _ := range args {
		argString += args[i]
		if i + 1 != len(args) {
			// this isn't the last arg, so add a space
			argString += " "
		}
	}

	command := exec.Command("go", "run", filepath, argString)
	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	command.Start()

	select {
		case <- done:
			return
	}
}

func watchForChange(filepath string, args []string) {
	done := make(chan bool, 1)
	initialFileInfo, _ := os.Stat(filepath)
	go runGoFile(filepath, args, done)
	for {
		time.Sleep(100 * time.Millisecond)
		fileInfo, _ := os.Stat(filepath)
		if initialFileInfo.ModTime() != fileInfo.ModTime() {
			initialFileInfo = fileInfo
			fmt.Printf("Change detected in %s, reloading...\n", filepath)
			close(done)
			done = make(chan bool, 1)
			go runGoFile(filepath, args, done)
		}
	}
}