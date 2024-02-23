package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func matchAndReplace(raw string, rule Rule) (string, bool, error) {
	r, err := regexp.Compile("(?i)" + rule.Match)
	if err != nil {
		return "", false, err
	}

	if !r.MatchString(raw) {
		return "", false, nil
	}

	items := r.FindAllStringSubmatch(raw, -1)
	matchItems := items[0][1:]

	replacerItems := []string{}
	for i, item := range matchItems {
		replacerItems = append(replacerItems, fmt.Sprintf("$%d", i+1))
		replacerItems = append(replacerItems, item)
	}

	return strings.NewReplacer(replacerItems...).Replace(rule.Rename), true, nil
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func moveFiles(from, to string, junks []string) error {
	files, err := os.ReadDir(from)
	if err != nil {
		return errors.Wrap(err, "read dir failed")
	}

	junks = append(junks, "metadata")
	for _, file := range files {
		baseName := file.Name()[:len(file.Name())-len(path.Ext(file.Name()))]

		match := false
		if strings.HasPrefix(baseName, "【") {
			toSmall, err := isFileToSmall(file)
			if err == nil && toSmall {
				match = true
			}
		}

		if !match {
			for _, junk := range junks {
				if baseName == junk {
					match = true
				}
			}
		}

		if match {
			continue
		}

		if err := os.Rename(filepath.Join(from, file.Name()), filepath.Join(to, file.Name())); err != nil {
			return err
		}
	}

	return nil
}

func isFileToSmall(f fs.DirEntry) (bool, error) {
	fi, err := f.Info()
	if err != nil {
		return false, err
	}

	return fi.Size() < 1024*1024, nil
}
