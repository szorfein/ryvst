package collection

import (
	"fmt"
	"os"
	"path/filepath"
)

func Scan(directory string) []string {
	var allTheYaml []string
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
	return allTheYaml
}
