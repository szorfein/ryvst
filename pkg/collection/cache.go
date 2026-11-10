// Library related to operations from ~/.cache/ryvst/
package collection

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ryvst/pkg/colour"
	"ryvst/pkg/unix"
)

// TODO: move some parts in internal/
var workdir = os.Getenv("HOME") + "/.cache/ryvst"

func CreateCacheDir(filename string) string {
	s := strings.Split(filename, ".yml")
	n := filepath.Join(workdir, s[0])
	unix.MkDAllP(n, 0o700)
	return n
}

func CompareTime(collection string, everytime time.Duration) bool {
	cacheDir := CreateCacheDir(collection)
	now := time.Now()

	lr := filepath.Join(cacheDir, "last.txt")
	if !unix.PathExist(lr) {
		return true
	}

	data, err := os.ReadFile(lr)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("Read %s: %d bytes\n", lr, len(data))
	date, err := time.Parse(time.RFC1123, string(data))
	if err != nil {
		fmt.Println("timeParse", err)
		os.Exit(1)
	}

	plus := date.Add(everytime)
	//fmt.Printf("date origin %s - modified with %s, %s\n", date, everytime, plus)
	if plus.After(now) {
		colour.Info(4, "Futur run: %s\n", plus.Format(time.RFC1123)) // not time, return false
	} else {
		return true
	}

	return false
}

func GenHashes(collection string, cachedir string, tasks []Tasks) {
	colour.Info(2, "Generating hashes for %s\n", collection)
	a := map[string]string{}
	for _, task := range tasks {
		if task.Name != "" {
			f := filepath.Join(cachedir, task.Name)
			a[task.Name] = unix.HashFile256(f)
		}
	}
	fmt.Println(a)
}

// Record Time if all the files has been downloaded for CompareTime().
func ApplyTimeNow(cacheDir string) {
	now := time.Now()
	p := filepath.Join(cacheDir, "last.txt")
	colour.Info(3, "Record Time.Now %s for %s\n", now.String(), cacheDir)

	file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	nowformat := now.Format(time.RFC1123)

	n, err := file.WriteString(string(nowformat))
	if err != nil {
		fmt.Println(n, err)
	}
	//fmt.Printf("wrote %d bytes\n", n)
}
