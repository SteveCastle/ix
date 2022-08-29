package main

import (
	"fmt"
	"os"
)

func findStore(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, f := range files {
		if f.Type().IsDir() && f.Name() == "index" {
			return path
		}
	}
	return findStore("../" + path)
}

func InitIndex() {
	err := os.Mkdir("index/", 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Tag(tag, file string) {
	store := findStore("./")
	fmt.Println("Creating tag assignment at:", store+"index/"+tag+"/"+file)
	err := os.Symlink(file, store+"index/"+tag+file)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateTag(tag string) {
	store := findStore("./")
	fmt.Println("Creating tag at:", store+"index/"+tag)
	err := os.MkdirAll(store+"index/"+tag, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Rebuild() {
	fmt.Println("Rebuilt index.")
}

func main() {
	args := os.Args[1:]
	cmd := args[0]

	switch cmd {
	case "init":
		InitIndex()
	case "create":
		CreateTag(args[1])
	case "tag":
		Tag(args[1], args[2])
	case "store":
		findStore("./")
	case "build":
		Rebuild()
	default:
		fmt.Println("Not a valid command.")
	}
}
