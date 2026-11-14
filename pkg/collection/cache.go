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

func getWorkdir() string { return os.Getenv("HOME") + "/.cache/ryvst" }

func CreateCacheDir(filename string) string {
	s := strings.Split(filename, ".yml")
	n := filepath.Join(getWorkdir(), s[0])
	unix.MkDAllP(n, 0o700)
	return n
}

// Return true if something would happen
func CompareTime(collection string, everytime time.Duration) (bool, error) {
	cacheDir := CreateCacheDir(collection)
	now := time.Now()

	lr := filepath.Join(cacheDir, "last.txt")
	if !unix.PathExist(lr) {
		return true, nil
	}

	data, err := os.ReadFile(lr)
	if err != nil {
		return false, fmt.Errorf("failed to read last time: %w", err)
	}

	date, err := time.Parse(time.RFC1123, string(data))
	if err != nil {
		return false, fmt.Errorf("failed to parse time: %w", err)
	}

	plus := date.Add(everytime)
	if !plus.After(now) { // yah, it's weird
		return true, nil
	}

	colour.Info(4, "Futur run: %s\n", plus.Format(time.RFC1123)) // not time, return false
	return false, nil
}

func GenHashes(cachedir string, files []string) error {
	colour.Info(2, "Generating hashes for %s\n", filepath.Base(cachedir))
	a := make(map[string]string)
	for _, file := range files {
		if file != "" {
			f := filepath.Join(cachedir, file)
			hash, err := unix.HashFile256(f)
			if err != nil {
				return err
			}
			a[file] = hash
		}
	}
	fmt.Println(a)
	return nil
}

// Record Time if all the files has been downloaded for CompareTime().
func ApplyTimeNow(cacheDir string) error {
	now := time.Now()
	p := filepath.Join(cacheDir, "last.txt")
	colour.Info(3, "Record Time.Now %s for %s\n", now.String(), cacheDir)

	file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	nowformat := now.Format(time.RFC1123)

	_, errws := file.WriteString(string(nowformat))
	if errws != nil {
		return err
	}
	return nil
}
