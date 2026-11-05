package main

import (
	"fmt"
	"os"

	"ryvst/internal/app"
	"ryvst/pkg/collection"
	"ryvst/pkg/colour"
	"ryvst/pkg/unix"
)

func main() {
	fmt.Printf("=== Ryvst %s ===\n", app.Version())
	config := app.Show_path()
	if config == "" {
		fmt.Println("No config found")
		os.Exit(1)
	}
	fmt.Printf("Configuraton located at %s\n", config)
	unix.MkdirP(config, 0o700)
	lists := collection.Scan(config)
	fmt.Println("Found:", len(lists), "collections")
	fmt.Println(lists)
	for _, list := range lists {
		collection.Read(config, list)
	}
	colour.RandPrintf(">")
	colour.RandPrintf(">")
	fmt.Println()
}
