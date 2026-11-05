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

func MkdirP(pathname string, permission os.FileMode) {
	if !PathExist(pathname) {
		fmt.Println("Creating directory", pathname)
		Mkdir(pathname, permission)
	}
}

func OpenFile(pathname string) *os.File {
	file, err := os.Open(pathname)
	if err != nil {
		fmt.Println("err collection.Read", err)
	}
	return file
}
