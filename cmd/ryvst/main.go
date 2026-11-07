package main

import (
	"fmt"
	"os"

	"ryvst/internal/app"
	"ryvst/pkg/collection"
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
	for _, list := range lists {
		c := collection.Read(config, list)
		if !collection.CompareTime(list, c.EveryTime) { // should return true or false, most logic?
			continue
		}
		fmt.Println("Continious flow...")
		//collection.Maker(c.Tasks)
	}
}
