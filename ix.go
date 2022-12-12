package ix

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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

func Tag(category, tag, filePath string) {
	store := FindStore("./")
	pwd, err := os.Getwd()
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
		basePath := path.Base(f)
		CreateTag(category, tag)
		dir := fmt.Sprintf("%s/%s/%s", store, category, tag)
		tagPath := path.Join(dir, basePath)
		filePath := path.Join(pwd, f)
		fmt.Fprintf(os.Stdout, "Tagging %s with %s/%s\n", filePath, category, tag)
		err = os.Link(filePath, tagPath)
		if err != nil {
			fmt.Println(err)
			return
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
