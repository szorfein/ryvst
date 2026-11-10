package download

import (
	"fmt"
	"os"
	"os/exec"

	"ryvst/pkg/colour"
	"ryvst/pkg/unix"
)

func gitClone(repoUrl string, dest string) error {
	colour.Info(3, "Cloning %s\n", repoUrl)
	_, err := exec.Command("git", "clone", repoUrl, dest).Output()
	if err != nil {
		return err
	}
	return nil
}

func gitPull(repoUrl string, dest string) error {
	err := os.Chdir(dest)
	if err != nil {
		return err
	}

	colour.Info(3, "Updating %s\n", repoUrl)
	_, err = exec.Command("git", "pull").Output()
	if err != nil {
		return err
	}
	return nil
}

func Git(repoUrl string, dest string) {
	if !unix.PathExist(dest) {
		err := gitClone(repoUrl, dest)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err := gitPull(repoUrl, dest)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
