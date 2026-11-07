package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"ryvst/pkg/colour"
)

func Grab(url string, dest string) {
	colour.Info("Downloading %s...\n", url)
	// TODO: http.Get don't have timeout
	// should we prefer resp, err := http.NewRequest(http.MethodGet, url, nil) ?
	// https://www.gofaq.org/en/how-to-read-an-http-response-body-in-go/
	resp, err := http.Get(url)
	if resp != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	size, _ := io.Copy(file, resp.Body)
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	colour.Info("%s | %d\n", dest, size)
}
