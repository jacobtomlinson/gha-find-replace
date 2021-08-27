package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"path/filepath"
	"regexp"
)

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
		includeRe := regexp.MustCompile(include)
		excludeRe := regexp.MustCompile(exclude)
		return includeRe.Match([]byte(path)) && !excludeRe.Match([]byte(path))
	}
	return false
}

func findAndReplace(path string, find string, replace string) (bool, error) {
	fmt.Println("looking for:")
	fmt.Println(find)
	if find != replace {
		read, readErr := ioutil.ReadFile(path)
		check(readErr)

		re := regexp.MustCompile(find)
		newContents := re.ReplaceAllString(string(read), replace)

		if newContents != string(read) {
			writeErr := ioutil.WriteFile(path, []byte(newContents), 0)
			check(writeErr)
			return true, nil
		}
	}

	return false, nil
}

func main() {
	include := os.Getenv("INPUT_INCLUDE")
	exclude := os.Getenv("INPUT_EXCLUDE")
	find := os.Getenv("INPUT_FIND")
	replace := os.Getenv("INPUT_REPLACE")

	reserved := []string{"INCLUDE","EXCLUDE","FIND","REPLACE"}

	INPUT_PREFIX := "INPUT_"

	files, filesErr := listFiles(include, exclude)
	check(filesErr)

	modifiedCount := 0

	if find!="" && replace!="" {
		for _, path := range files {
			modified, findAndReplaceErr := findAndReplace(path, find, replace)
			check(findAndReplaceErr)
	
			if modified {
				modifiedCount++
			}
		}
	}

	for _, pair := range os.Environ() {
		if strings.Contains(pair,INPUT_PREFIX) {
			keyValue := strings.SplitN(pair,"=",2)
			find := strings.SplitN(keyValue[0],"_",2)[1]
			replace := keyValue[1]

			i := sort.Search(len(reserved), func(i int) bool { return reserved[i] == find })

			if i == len(reserved) {
				files, filesErr := listFiles(include, exclude)
				check(filesErr)
				for _, path := range files {
					modified, findAndReplaceErr := findAndReplace(path, "(?i)"+find, replace)
					check(findAndReplaceErr)
		
					if modified {
						modifiedCount++
					}
				}
			}
					
		}
	}

	fmt.Println(fmt.Sprintf(`::set-output name=modifiedFiles::%d`, modifiedCount))
}
