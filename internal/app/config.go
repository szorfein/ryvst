package app

import (
	"fmt"
	"os"
)

func Show_path() string {
	config := "" // default void
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		config = os.Getenv("XDG_CONFIG_HOME") + "/ryvst"
		fmt.Printf("xdg_config_home exist %s\n", config)
	} else if os.Getenv("HOME") != "" {
		config = os.Getenv("HOME") + "/.config/ryvst"
		fmt.Printf("else Home exist %s\n", config)
	}
	return config
}
