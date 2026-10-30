package unix

import (
	"fmt"
	"os"
)

func PathExist(pathname string) bool {
	_, err := os.Stat(pathname)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

func Mkdir(pathname string, permission os.FileMode) {
	if err := os.Mkdir(pathname, permission); err != nil {
		if os.IsExist(err) {
			fmt.Println("Directory already exists.")
		} else {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Directory created successfully.")
	}
}
