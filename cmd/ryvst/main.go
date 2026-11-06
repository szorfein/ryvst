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
	cacheDir := os.Getenv("HOME") + "/.cache/ryvst"
	unix.MkdirP(config, 0o700)
	unix.MkdirP(cacheDir, 0o700)
	lists := collection.Scan(config)
	for _, list := range lists {
		c := collection.Read(config, list)
		collectionCacheDir := unix.CreateCacheDir(cacheDir, list)
		colour.Info("Using cache %s\n", collectionCacheDir)
		collection.CompareTime(collectionCacheDir, c.EveryTime)
	}
}
