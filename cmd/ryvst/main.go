package main

import (
	"fmt"

	"ryvst/version"
)

func main() {
	v := version.Display()
	fmt.Printf("=== Ryvst %s ===\n", v)
}
