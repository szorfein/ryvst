package collection

import (
	"path/filepath"

	"ryvst/pkg/download"
)

func Maker(collection string, tasks []Tasks) {
	c := CreateCacheDir(collection)

	// TODO
	// 1. If a SIGNATURE file is set and cacheDir/hashes.yml exist, download the sig file and check if need to download
	// 2. Downloads or Skip the file
	// 3. Generate hash and save them to cacheDir/hashes.yml
	// 4. Update cacheDir/last-run.yml
	// 5. Move or Extract to the right path Tasks.To_dest depend of this type file (rar, zip, ttf...)
	for _, task := range tasks {
		f := filepath.Join(c, task.Name)

		download.Grab(task.Url, f) // 2, TODO with goroutine
	}
}
