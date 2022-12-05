package ix

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func findStore(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, f := range files {
		if f.Type().IsDir() && f.Name() == "ix" {
			return path
		}
	}
	return findStore("../" + path)
}

func InitIndex() {
	err := os.Mkdir("ix/", 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Tag(category, tag, filePath string) {
	store := findStore("./")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	files := []string{}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println("Creating tag assignment at:", store+"ix/"+tag+"/"+f)
		CreateTag(category, tag)
		dir := fmt.Sprintf("%s/ix/%s/%s", store, category, tag)
		err = os.Link(path.Join(pwd, f), path.Join(dir, basePath))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func CreateTag(category, tag string) {
	store := findStore("./")
	path := fmt.Sprintf("%s/ix/%s/%s", store, category, tag)
	fmt.Println("creating tag with category:", category, "tag:", tag)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func CrossIndex(parent, child string) {
	store := findStore("./")
	parentItems, err := os.ReadDir(store + "ix/" + parent)
	if err != nil {
		fmt.Println(err)
		return
	}
	childItems, err := os.ReadDir(store + "ix/" + child)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range parentItems {
		fmt.Println("Parent Items:", p.Name())
	}
	for _, c := range childItems {
		fmt.Println("Child Items:", c.Name())
	}
}

func Rebuild() {
	fmt.Println("Rebuilt index.")
}
