package unix

import (
	"fmt"
	"os"
)

func PathExist(pathname string) bool {
	_, err := os.Stat(pathname)
	return err == nil
}

func Mkdir(pathname string, permission os.FileMode) {
	if err := os.Mkdir(pathname, permission); err != nil {
		if !os.IsExist(err) {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Directory created successfully.")
	}
}
