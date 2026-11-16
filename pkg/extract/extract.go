// helpers for the extract package
package extract

import (
	"path/filepath"
	"strings"
)

// Remove first directory if found suffix -main or -master
func stripComponent(localname string, strip bool) string {
	normalizedName := filepath.ToSlash(localname)
	//split := strings.Split(localname, string(os.PathSeparator))
	split := strings.SplitN(normalizedName, "/", 2)

	if len(split) != 2 {
		return localname
	}
	if strip {
		return filepath.FromSlash(split[1])
	}
	if strings.HasSuffix(split[0], "-main") {
		return filepath.FromSlash(split[1])
	}
	if strings.HasSuffix(split[0], "-master") {
		return filepath.FromSlash(split[1])
	}

	return localname
}
