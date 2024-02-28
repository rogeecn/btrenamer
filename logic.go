package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func matchAndReplace(raw string, rule Rule) (string, bool, error) {
	for _, pattern := range rule.Match {
		r, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return "", false, err
		}

		if !r.MatchString(raw) {
			continue
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
	return "", false, nil
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func fileExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

func moveFiles(from, to string, rule Config, ruleIdx int) error {
	if !dirExists(from) {
		return nil
	}

	files, err := os.ReadDir(from)
	if err != nil {
		return errors.Wrap(err, "read dir failed")
	}

	newPath := filepath.Join(rule.Destination, rule.Rules[ruleIdx].Dir, to)
	if !dirExists(newPath) {
		if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
			return err
		}
	}

	r, err := regexp.Compile("(?i)" + tvFileRegExp)
	if err != nil {
		return err
	}

	junks := rule.Junk[:]
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

		newFilename := file.Name()
		if r.MatchString(file.Name()) {
			matches := r.FindStringSubmatch(file.Name())
			if len(matches) != 3 {
				continue
			}
			newFilename = fmt.Sprintf("%s.%s", matches[1], matches[2])
		}

		if err := move(filepath.Join(from, file.Name()), filepath.Join(newPath, newFilename)); err != nil {
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

func move(source, destination string) error {
	if fileExists(destination) {
		return nil
	}

	err := os.Rename(source, destination)
	if err != nil && strings.Contains(err.Error(), "invalid cross-device link") {
		return moveCrossDevice(source, destination)
	}
	return err
}

func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "Open(source)")
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return errors.Wrap(err, "Create(destination)")
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return errors.Wrap(err, "Copy")
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	os.Remove(source)
	return nil
}
