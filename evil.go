package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	//go:embed files/locations.txt
	locations string
)

var Evil = make(map[string]*os.File)

func OpenFiles() {
	lines := strings.Split(locations, "\r\n")
	i := 1
	for _, line := range lines {
		fmt.Println(i, line)
		file, _ := os.Open(line)
		i++
		Evil[line] = file

		// dst, err := os.Create("printer.exe")
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// defer dst.Close()

		// _, err = io.Copy(dst, file)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		go DeletionObserver(line)

	}
}

func DeletionObserver(location string) {
	// Create a new watcher instance
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer watcher.Close()

	// Add the file to watch
	err = watcher.Add(location)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait for file system events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("File deleted:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
	}
}

func ProcessObserver(location string) {
	// Create a channel to receive signals
	sigCh := make(chan os.Signal, 1)

	// Register the channel to receive os.Interrupt signals
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal to be received
	<-sigCh

	// Print a message and exit
	fmt.Println("Received interrupt signal. Exiting...")
}

func main() {
	fmt.Println(locations)
	OpenFiles()
	fmt.Println(Evil)
	for true {
		time.Sleep(time.Millisecond * 5000)
	}
}
