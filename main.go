package main

import (
	"fmt"
	"os"
	"path"
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

func Tag(tag, file string) {
	store := findStore("./")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Creating tag assignment at:", store+"ix/"+tag+"/"+file)
	CreateTag(tag)
	err = os.Link(path.Join(pwd, file), path.Join(store+"ix/"+tag, file))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateTag(tag string) {
	store := findStore("./")
	path := store + "ix/" + tag
	fmt.Println("Creating tag at:", path)
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("IX v0.01")
		return
	}
	args := os.Args[1:]
	cmd := args[0]

	switch cmd {
	case "init":
		InitIndex()
	case "create":
		CreateTag(args[1])
	case "tag":
		Tag(args[1], args[2])
	case "cross":
		CrossIndex(args[1], args[2])
	case "store":
		findStore("./")
	case "build":
		Rebuild()
	default:
		fmt.Println("Not a valid command.")
	}
}
