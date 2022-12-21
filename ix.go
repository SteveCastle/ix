package ix

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/segmentio/ksuid"
)

var fileTypes = []string{".jpg", ".jpeg", ".gif", ".png", ".mp4"}

type Config struct {
	Store string
}

func saveConfig(file string, config *Config) error {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(config)
}

func loadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func FindStore(path string) string {
	pathVolume := filepath.VolumeName(path)
	var storePath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if pathVolume == filepath.VolumeName(homeDir) {
		storePath = homeDir + "/.ix"
	} else {
		storePath = pathVolume + "/.ix"
	}
	config, err := loadConfig(storePath + "/config.json")
	if err != nil {
		defaultConfig := &Config{Store: storePath}
		saveConfig(storePath+"/config.json", defaultConfig)
		return defaultConfig.Store
	}
	return config.Store
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
		// Return the entire input string.
		return input
	}

	// Return the part of the input string after the final slash character.
	return input[lastSlashIndex+1:]
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Tag(category, tag, inputPath string) {
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		fmt.Println("Error getting absolute path", err)
	}
	store := FindStore(absolutePath)
	tagDirectory := fmt.Sprintf("%s/%s/%s", store, category, tag)
	// If tag directory does not exist, then create it.
	if _, err := os.Stat(tagDirectory); os.IsNotExist(err) {
		CreateTag(category, tag, absolutePath)
	}

	files := []string{}
	fileInfo, err := os.Stat(absolutePath)
	if err != nil {
		fmt.Println("Error opening file", err)
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(absolutePath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				ext := filepath.Ext(path)
				if !info.IsDir() && contains(fileTypes, ext) {
					files = append(files, path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	} else {
		files = append(files, absolutePath)
	}
	for _, f := range files {
		fileName := afterLastSlash(f)
		fileBase := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		fileExt := filepath.Ext(fileName)
		newFileName := fmt.Sprintf("%s-%s%s", fileBase, ksuid.New(), fileExt)
		tagPath := path.Join(tagDirectory, newFileName)
		filePath := f
		fmt.Fprintf(os.Stdout, "Tagging %s with %s/%s\n", filePath, category, tag)
		err = os.Link(filePath, tagPath)
		if err != nil {
			fmt.Println("Could not create link: ", err)
		}
	}
}

func CreateTag(category, tag string, rootPath string) {
	store := FindStore(rootPath)
	path := fmt.Sprintf("%s/%s/%s", store, category, tag)
	fmt.Println("creating tag at path:", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func ListCategories() []string {
	store := FindStore("./")
	categories := []string{}
	files, err := os.ReadDir(store)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			categories = append(categories, f.Name())
		}
	}
	return categories
}

func ListTags(category string) []string {
	store := FindStore("./")
	tags := []string{}
	path := fmt.Sprintf("%s/%s", store, category)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			tags = append(tags, f.Name())
		}
	}
	return tags
}
