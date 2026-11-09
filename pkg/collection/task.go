package collection

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"ryvst/pkg/download"
)

func PushToList(collection string, tasks []Tasks) []map[string]string {
	var arr []map[string]string
	for _, task := range tasks {
		switch filepath.Ext(task.Url) {
		case ".zip":
			amap := map[string]string{"Name": task.Name, "Url": task.Url}
			arr = append(arr, amap)
		default:
			fmt.Println("None yet supported:", task.Url)
		}
	}
	return arr
}

func DownloadAllFiles(collection string, tasks []map[string]string) {
	cache := CreateCacheDir(collection)

	var waitGroup sync.WaitGroup
	for _, task := range tasks {
		waitGroup.Go(func() {
			f := filepath.Join(cache, task["Name"])
			download.Grab(task["Url"], f)
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
	GenHashes(collection, c, tasks)     // 3, TODO Save them, Maybe in YAML or JSON ?
	err := ToDest(collection, c, tasks) // 4
	if err != nil {
		log.Fatal(err)
	}
	ApplyTimeNow(c) // 5
}

func DownloadAll(config string, lists []string) {
	downloadsList := make(map[string][]map[string]string) // don't work withou make()
	for _, list := range lists {
		c := Read(config, list)
		if !CompareTime(list, c.EveryTime) { // should return true or false, most logic?
			continue
		}
		downloadsList[list] = PushToList(list, c.Tasks)
	}
	for i, d := range downloadsList {
		DownloadAllFiles(i, d)
	}
}
