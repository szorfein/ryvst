// Library related to operations from ~/.cache/ryvst/
package collection

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ryvst/pkg/unix"
)

// TODO: move some parts in internal/
var workdir = os.Getenv("HOME") + "/.cache/ryvst"

func CreateCacheDir(cachedir string, filename string) string {
	s := strings.Split(filename, ".yml")
	n := filepath.Join(cachedir, s[0])
	unix.MkDAllP(n, 0o700)
	return n
}

func CompareTime(collection string, everytime time.Duration) bool {
	cacheDir := CreateCacheDir(workdir, collection)

	// TODO
	// 1. Should load cacheDir/last-run.yml, return false if no exist
	// 2. If exist, read the value "lastTime" and compare to everyTime
	fmt.Println(everytime)
	lr := filepath.Join(cacheDir, "last-run.yml")
	if !unix.PathExist(lr) {
		fmt.Println("lastrun doesn't exist", lr)
		return true
	}

	return false
}
