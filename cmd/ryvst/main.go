package main

import (
	"fmt"
	"os"

	"ryvst/internal/app"
	"ryvst/pkg/collection"
	"ryvst/pkg/unix"
)

func main() {
	v := app.Version()
	fmt.Printf("=== Ryvst %s ===\n", v)
	config := app.Show_path()
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
	for _, list := range lists {
		collection.Read(config, list)
	}
}
