package collection

import (
	"fmt"
	"os"
	"path/filepath"

	"ryvst/pkg/colour"
)

func Scan(directory string) []string {
	var allTheYaml []string
	colour.Info(2, "Scanning collections from %s\n", directory)

	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		match, _ := filepath.Match("*.yml", d.Name())
		if match { // true if match
			allTheYaml = append(allTheYaml, d.Name())
		}
		return nil
	})
	if err != nil {
		fmt.Println("Got Error:", err)
	}

	colour.Info(3, "Found %d collections\n", len(allTheYaml))
	return allTheYaml
}
