package collection

import (
	"fmt"
	"time"
)

func CompareTime(cacheDir string, everytime time.Duration) bool {
	// TODO
	// 1. Should load cacheDir/metadata.yaml, return false if no exist
	// 2. If exist, read the value "lastTime" and compare to everyTime
	fmt.Println(cacheDir)
	fmt.Println(everytime)
	return false
}
