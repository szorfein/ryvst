package collection

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

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
	Execute   string        `yaml:"execute"`
	Tasks     []Tasks
}

func Read(dir string, filename string) Config {
	fullpath := filepath.Join(dir, filename)
	file := unix.OpenFile(fullpath)

	// Add this safety check
	if file == nil {
		fmt.Println("Error: file is nil, cannot read config")
		os.Exit(1)
	}

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
		Execute:   "",
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Err collection.Unmarshal", err)
		os.Exit(1)
	}
	//fmt.Printf("YAML content: %+v\n", config)
	return config
}
