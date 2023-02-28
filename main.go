package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var news []string
	dir := "news"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			path := filepath.Join(dir, file.Name())
			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			news = append(news, string(data))
		}
	}
	fmt.Println(news[0])
}
