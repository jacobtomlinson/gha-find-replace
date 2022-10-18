package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gobwas/glob"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func getenvInt(key string, def int) (int, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func getenvBool(key string) (bool, error) {
	s, err := getenvStr(key)
	if err != nil {
		return true, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return true, err
	}
	return v, nil
}

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func listFiles(include string, exclude string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if doesFileMatch(path, include, exclude) {
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList, err
}

func doesFileMatch(path string, include string, exclude string) bool {
	if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
		includeGlob := glob.MustCompile(include)
		excludeGlob := glob.MustCompile(exclude)
		return includeGlob.Match(path) && !excludeGlob.Match(path)
	}
	return false
}

func findAndReplace(path string, find string, replace string, regex bool) (bool, error) {
	if find != replace {
		read, readErr := ioutil.ReadFile(path)
		check(readErr)

		var newContents = ""
		if regex {
			re := regexp.MustCompile(find)
			newContents = re.ReplaceAllString(string(read), replace)
		} else {
			newContents = strings.ReplaceAll(string(read), find, replace)
		}

		if newContents != string(read) {
			writeErr := ioutil.WriteFile(path, []byte(newContents), 0)
			check(writeErr)
			return true, nil
		}
	}

	return false, nil
}

func setGithubEnvOutput(key string, value int) {
	outputFilename := os.Getenv("GITHUB_ENV")
	f, err := os.OpenFile(outputFilename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(fmt.Sprintf(`%s=%d`, key,  value)); err != nil {
			log.Println(err)
		}
}

func main() {
	include, _ := getenvStr("INPUT_INCLUDE")
	exclude, _ := getenvStr("INPUT_EXCLUDE")
	find, findErr := getenvStr("INPUT_FIND")
	replace, replaceErr := getenvStr("INPUT_REPLACE")
	regex, regexErr := getenvBool("INPUT_REGEX")

	if findErr != nil {
		panic(errors.New("gha-find-replace: expected with.find to be a string"))
	}

	if replaceErr != nil {
		panic(errors.New("gha-find-replace: expected with.replace to be a string"))
	}

	if regexErr != nil {
		regex = true
	}

	files, filesErr := listFiles(include, exclude)
	check(filesErr)

	modifiedCount := 0

	for _, path := range files {
		modified, findAndReplaceErr := findAndReplace(path, find, replace, regex)
		check(findAndReplaceErr)

		if modified {
			modifiedCount++
		}
	}

	setGithubEnvOutput("modifiedFiles", modifiedCount)
}
