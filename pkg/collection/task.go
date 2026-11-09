package collection

import (
	"log"
	"path/filepath"
	"sync"

	"ryvst/pkg/download"
)

func DownloadAllFiles(cache string, tasks []Tasks) {
	var waitGroup sync.WaitGroup
	for _, task := range tasks {
		waitGroup.Go(func() {
			f := filepath.Join(cache, task.Name)
			download.Grab(task.Url, f)
		})
	}
	waitGroup.Wait()
}

func Maker(collection string, tasks []Tasks) {
	c := CreateCacheDir(collection)

	// TODO
	// 1. If a SIGNATURE file is set and cacheDir/hashes.yml exist, download the sig file and check if need to download
	// 2. Downloads or Skip the file
	// 3. Generate hash and save them to cacheDir/hashes.yml
	// 4. Move or Extract to the right path Tasks.To_dest depend of this type file (rar, zip, ttf...)
	// 5. Update cacheDir/last-run.yml
	DownloadAllFiles(c, tasks)          // 2
	GenHashes(collection, c, tasks)     // 3, TODO Save them, Maybe in YAML or JSON ?
	err := ToDest(collection, c, tasks) // 4
	if err != nil {
		log.Fatal(err)
	}
	ApplyTimeNow(c) // 5
}
