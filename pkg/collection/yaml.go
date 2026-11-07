package collection

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"ryvst/pkg/colour"
	"ryvst/pkg/unix"

	"go.yaml.in/yaml/v4"
)

type Tasks struct {
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	ToDest string `yaml:"to_dest"`
}

type Config struct {
	EveryTime time.Duration `yaml:"every_time"`
	Tasks     []Tasks
}

func Read(dir string, filename string) Config {
	fullpath := filepath.Join(dir, filename)
	colour.Info("Loading collection %s\n", fullpath)
	file := unix.OpenFile(fullpath)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file", err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("err collection.ReadAll", err)
	}
	config := Config{ // Setting default values
		EveryTime: 24 * time.Hour,
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("err collection.Unmarshal", err)
	}
	fmt.Printf("YAML content: %+v\n", config)
	return config
}
