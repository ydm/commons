package commons

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func checkFile(file fs.FileInfo, prefix, suffix string) bool {
	if file.IsDir() {
		return false
	}

	if prefix != "" && !strings.HasPrefix(file.Name(), prefix) {
		return false
	}

	if suffix != "" && !strings.HasSuffix(file.Name(), suffix) {
		return false
	}

	return true
}

func ListFiles(input, prefix, suffix string) []string {
	// Create result.
	ans := make([]string, 0, 16)

	// If given parameter is actually a single file, check it and return.
	file, err := os.Stat(input)
	if err != nil {
		Msg(log.Warn().Err(err))

		return ans
	}

	if !file.IsDir() && checkFile(file, prefix, suffix) {
		ans = append(ans, input)

		return ans
	}

	// Iterate over all directories recursively and filter files of interest.

	dirs := []string{input}
	for len(dirs) > 0 {
		head := dirs[0]
		dirs = dirs[1:]

		files, err := ioutil.ReadDir(head)
		if err != nil {
			Msg(log.Warn().Err(err))

			continue
		}

		for _, file := range files {
			path := filepath.Join(head, file.Name())

			if file.IsDir() {
				dirs = append(dirs, path)
			}

			if checkFile(file, prefix, suffix) {
				ans = append(ans, path)
			}
		}
	}

	return ans
}
