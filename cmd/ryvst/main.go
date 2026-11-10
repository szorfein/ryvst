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
	unix.MkdirP(config, 0o700)
	lists := collection.Scan(config)
	// First download all
	done, doneGit := collection.DownloadAll(config, lists)
	if done == 0 && doneGit == 0 {
		colour.Info(2, "Bye\n")
		os.Exit(0)
	}
	// Next hashes, move files and apply Time now
	for _, list := range lists {
		c := collection.Read(config, list)
		if !collection.CompareTime(list, c.EveryTime) { // should return true or false, most logic?
			continue
		}
		collection.Maker(list, c.Tasks, done)
	}
}
