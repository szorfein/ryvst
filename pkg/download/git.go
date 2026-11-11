package download

import (
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
	colour.Info(3, "Updating %s\n", repoUrl)
	cmd := exec.Command("git", "pull")
	cmd.Dir = dest
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func Git(repoUrl string, dest string) error {
	if !unix.PathExist(dest) {
		return gitClone(repoUrl, dest)
	}
	return gitPull(repoUrl, dest)
}
