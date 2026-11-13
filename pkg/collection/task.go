package collection

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"ryvst/pkg/colour"
	"ryvst/pkg/download"
)

// Push non-Git url to the list
// For now, only zip file works but would accept more (jpg, png), (pdf, epub), tar.gz, (otf, ttf)
func PushToHTTPSList(collection string, tasks []Tasks) []map[string]string {
	var arr []map[string]string
	for _, task := range tasks {
		mimetype, err := download.GetFileContentType(task.Url)
		if err != nil {
			fmt.Println("Error detect type", err)
			os.Exit(1)
		}
		switch mimetype {
		case "application/zip", "application/x-gzip":
			amap := map[string]string{"Name": task.Name, "Url": task.Url}
			arr = append(arr, amap)
		case "font/ttf", "text/xml; charset=utf-8":
			amap := map[string]string{"Name": task.Name, "Url": task.Url}
			arr = append(arr, amap)
		case "application/octet-stream": // octet-stream is returned when the mimetype is no found...
			if strings.HasSuffix(task.Name, ".tar.xz") {
				amap := map[string]string{"Name": task.Name, "Url": task.Url}
				arr = append(arr, amap)
			}
		default:
			if !strings.HasSuffix(task.Url, ".git") {
				fmt.Println("None yet supported:", mimetype)
			}
		}
	}
	return arr
}

// Push Git repository to list
// Url should end by .git, e.g: https://github.com/szorfein/ryvst.git
func PushToGitList(collection string, tasks []Tasks) []map[string]string {
	var arr []map[string]string
	for _, task := range tasks {
		switch filepath.Ext(task.Url) {
		case ".git":
			amap := map[string]string{"ToDest": task.ToDest, "Url": task.Url}
			arr = append(arr, amap)
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
			err := download.Https(task["Url"], f)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		})
	}
	waitGroup.Wait()
}

func gitAllRepos(collection string, tasks []map[string]string) {
	home := os.Getenv("HOME")
	var waitGroup sync.WaitGroup
	for _, task := range tasks {
		waitGroup.Go(func() {
			f := filepath.Join(home, task["ToDest"])
			err := download.Git(task["Url"], f)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		})
	}
	waitGroup.Wait()
}

func Maker(collection string, tasks []Tasks, done int) {
	c := CreateCacheDir(collection)

	// TODO
	// 1. If a SIGNATURE file is set and cacheDir/hashes.yml exist, download the sig file and check if need to download
	// 2. Downloads or Skip the file
	// 3. Generate hash and save them to cacheDir/hashes.yml
	// 4. Move or Extract to the right path Tasks.To_dest depend of this type file (rar, zip, ttf...)
	// 5. Update cacheDir/last-run.yml
	if done >= 1 {
		GenHashes(collection, c, tasks)     // 3, TODO Save them, Maybe in YAML or JSON ?
		err := ToDest(collection, c, tasks) // 4
		if err != nil {
			log.Fatal(err)
		}
	}
	ApplyTimeNow(c) // 5
}

func DownloadAll(config string, lists []string) (int, int) {
	downloadsList := make(map[string][]map[string]string) // don't work withou make()
	gitList := make(map[string][]map[string]string)       // don't work withou make()
	taskDone, taskGitDone := 0, 0
	for _, list := range lists {
		c := Read(config, list)
		if !CompareTime(list, c.EveryTime) { // should return true or false, most logic?
			continue
		}
		downloadsList[list] = PushToHTTPSList(list, c.Tasks)
		gitList[list] = PushToGitList(list, c.Tasks)
	}
	for i, d := range downloadsList {
		if len(d) == 0 {
			continue
		}
		taskDone++
		DownloadAllFiles(i, d)
	}

	// Git - clone or pull
	for i, d := range gitList {
		if len(d) == 0 {
			continue
		}
		taskGitDone++
		gitAllRepos(i, d)
	}

	return taskDone, taskGitDone
}

func RunCommand(cmd string) error {
	if cmd == "" {
		return nil
	}
	colour.Info(2, "Executing command %s\n", cmd)
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
