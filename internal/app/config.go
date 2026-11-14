package app

import (
	"os"

	"ryvst/pkg/unix"
)

func Show_path() string {
	config := "" // default void
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		config = os.Getenv("XDG_CONFIG_HOME") + "/ryvst"
	} else if os.Getenv("HOME") != "" {
		config = os.Getenv("HOME") + "/.config/ryvst"
	}
	unix.MkdirP(config, 0o700)

	return config
}
