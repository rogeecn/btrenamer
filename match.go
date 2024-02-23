package main

import (
	"fmt"
	"os"
	"regexp"
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
	matchItems := []any{}
	for _, item := range items[0][1:] {
		matchItems = append(matchItems, any(item))
	}
	return fmt.Sprintf(rule.Rename, matchItems...), true, nil
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
