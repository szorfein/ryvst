package collection

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func RunCommand(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
