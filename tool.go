package main

import (
	"os"
	"regexp"
)

// find seq pad return. ex) .1001 or .%04d
func FindSeqPad(str string) string {
	re, _ := regexp.Compile(".+([_.][0-9]+)(.[a-zA-Z]+$)")
	results := re.FindStringSubmatch(str)
	if results != nil {
		return results[1]
	}
	re, _ = regexp.Compile(".+([_.]%[0-9]+d)(.[a-zA-Z]+$)")
	results = re.FindStringSubmatch(str)
	if results != nil {
		return results[1]
	}
	return ""
}

// check file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
