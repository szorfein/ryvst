package main

import (
	"fmt"
	"os"

	"ryvst/internal/app"
	"ryvst/pkg/collection"
	"ryvst/pkg/unix"
	"ryvst/version"
)

func main() {
	v := version.Display()
	fmt.Printf("=== Ryvst %s ===\n", v)
	config := app.Show_path()
	fmt.Println(config)
	if config == "" {
		fmt.Println("No config found")
		os.Exit(1)
	}
	fmt.Printf("Configuraton located at %s\n", config)
	if !unix.PathExist(config) {
		fmt.Println("Creating directory", config)
		unix.Mkdir(config, 0o700)
	}
	lists := collection.Scan(config)
	fmt.Println("Found:", len(lists), "collections")
	fmt.Println(lists)
}
