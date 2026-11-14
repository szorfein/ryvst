package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"ryvst/internal/app"
	"ryvst/pkg/collection"
	"ryvst/pkg/colour"
	"ryvst/pkg/download"
)

type Download struct {
	Name      string
	Url       string
	Cache     string
	ToDest    string
	EveryTime time.Duration
	IsGit     bool
	Command   string
}

// if we want to Collect the Mimetype here, we would use goroutine...
func loopInTask(d *[]Download, list string, config string) {
	c := collection.Read(config, list)
	cache := collection.CreateCacheDir(list)

	if !collection.CompareTime(list, c.EveryTime) { // should return true or false, most logic?
		return
	}

	for _, task := range c.Tasks {
		todest := os.Getenv("HOME") + "/" + task.ToDest
		isGit := false
		if strings.HasSuffix(task.Url, ".git") {
			isGit = true
		}
		*d = append(*d, Download{Name: task.Name, Url: task.Url, Cache: cache, ToDest: todest, EveryTime: c.EveryTime, IsGit: isGit, Command: c.Execute})
	}
}

func myPush(d *[]Download, lists []string, config string) {
	colour.Info(2, "Amassing tasks from...")
	for _, list := range lists {
		fmt.Printf(" %s", list)
		loopInTask(d, list, config)
	}
	fmt.Printf(" [done]\n")
}

func getGit(url string, dest string) {
	f := dest
	err := download.Git(url, f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getHttps(url string, cache string, name string) {
	f := filepath.Join(cache, name)
	err := download.Https(url, f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// TODO Limit thread here?
func downloadAllFiles(d []Download) {
	var waitGroup sync.WaitGroup
	for _, v := range d {
		waitGroup.Go(func() {
			if v.IsGit {
				getGit(v.Url, v.ToDest)
			} else {
				getHttps(v.Url, v.Cache, v.Name)
			}
		})
	}
	waitGroup.Wait()
}

func movingFiles(d []Download) {
	colour.Info(2, "Moving files...\n")
	for _, v := range d {
		if v.IsGit {
			continue
		}
		err := collection.ToDest(v.Cache, v.Name, v.ToDest)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func executeCommands(d []Download) {
	uniqueCacheDir := make(map[string]string)
	for _, v := range d {
		uniqueCacheDir[v.Cache] = v.Command // Build unique value cacheDir/command
	}
	for k, v := range uniqueCacheDir { // k=key, v=command
		if v == "" {
			continue
		}
		colour.Info(4, "Execute %s for %s...\n", v, k)
		err := collection.RunCommand(v)
		if err != nil {
			fmt.Println("Err with command", err)
			os.Exit(1)
		}
	}
}

func recordMetas(d []Download) {
	uniqueCacheDir := make(map[string][]string)
	for _, v := range d {
		cache := v.Cache
		if cache == v.Cache {
			uniqueCacheDir[cache] = append(uniqueCacheDir[cache], v.Name) // Build unique value cacheDir
		}
	}
	for k := range uniqueCacheDir { // k=key, v=command
		collection.GenHashes(k, uniqueCacheDir[k])
		collection.ApplyTimeNow(k)
	}
}

func main() {
	fmt.Printf("=== Ryvst %s ===\n", app.Version())
	config := app.Show_path()
	if config == "" {
		fmt.Println("No config found")
		os.Exit(1)
	}
	lists := collection.Scan(config)

	var s []Download
	k := &s // pointer on []Download

	myPush(k, lists, config)

	downloadAllFiles(s)
	movingFiles(s)
	executeCommands(s)
	recordMetas(s)

	os.Exit(0)
}
