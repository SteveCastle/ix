package ix

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FindStore(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, f := range files {
		if f.Type().IsDir() && f.Name() == "ix" {
			return path + "ix"
		}
	}
	return FindStore("../" + path)
}

func InitIndex() {
	err := os.Mkdir("ix/", 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Returns the part of the input string after the final slash character.
// If there is no slash character in the input string, returns the empty string.
func afterLastSlash(input string) string {
	// Find the index of the final slash character.
	// If there is no slash character in the input string,
	// strings.LastIndex returns -1.
	lastSlashIndex := strings.LastIndex(input, "\\")

	if lastSlashIndex == -1 {
		// There is no slash character in the input string.
		// Return the empty string.
		return ""
	}

	// Return the part of the input string after the final slash character.
	return input[lastSlashIndex+1:]
}

func Tag(category, tag, filePath string) {
	store := FindStore("./")
	pwd, err := os.Getwd()
	tagDirectory := fmt.Sprintf("%s/%s/%s", store, category, tag)
	// If tag directory does not exist, then create it.
	if _, err := os.Stat(tagDirectory); os.IsNotExist(err) {
		CreateTag(category, tag)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	files := []string{}
	fileInfo, err := os.Stat(path.Join(pwd, filePath))
	if err != nil {
		fmt.Println("Error opening file", err)
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(filePath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	} else {
		files = append(files, filePath)
	}
	for _, f := range files {
		fileName := afterLastSlash(f)
		//Console log the file name
		fmt.Println("F NAME:", f)
		tagPath := path.Join(tagDirectory, fileName)
		filePath := path.Join(pwd, f)
		fmt.Fprintf(os.Stdout, "Tagging %s with %s/%s\n", filePath, category, tag)
		err = os.Link(filePath, tagPath)
		if err != nil {
			fmt.Println("Could not create link: ", err)
		}
	}
}

func CreateTag(category, tag string) {
	store := FindStore("./")
	path := fmt.Sprintf("%s/%s/%s", store, category, tag)
	fmt.Println("creating tag with category:", category, "tag:", tag)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
