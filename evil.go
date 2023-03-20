package main

import (
	_ "embed"
	"fmt"
)

var (
	//go:embed files/locations.txt
	locations string
)

func main() {
	fmt.Println(locations)
}
