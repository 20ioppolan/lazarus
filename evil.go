package main

import (
	"embed"
	_ "embed"
	"fmt"
)

var (
	//go:embed files
	content embed.FS
)

func main() {
	entries, err := content.ReadDir("files")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Println(entry.Name())
		}
	}
}
