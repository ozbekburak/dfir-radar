package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/completion"
)

var (
	promptPrefix = "I will give you an articles about DFIR in various categories Please extract maximum three keywords from the articles to create insights. Use the following format: Keywords: keyword1, keyword2, keyword3. Use only lowercase letters \n"
	newsDir      = "news"
	textList     []string
)

func main() {
	newsFiles, err := os.ReadDir(newsDir)
	if err != nil {
		log.Printf("Failed to read news directory '%s' error: %v", newsDir, err)
		os.Exit(1)
	}
	for _, file := range newsFiles {
		if !file.IsDir() && file.Type().IsRegular() {
			path := filepath.Join(newsDir, file.Name())
			text, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Failed to read file '%s' error: %v", path, err)
				continue
			}
			textList = append(textList, string(text))
		}
	}

	extractKeywords(textList)
}

func extractKeywords(texts []string) {
	ctx := context.Background()
	s := openai.NewSession(os.Getenv("OPENAI_API_KEY"))
	client := completion.NewClient(s, "text-davinci-003")
	var keywordList []string
	for _, text := range texts {
		resp, err := client.Create(ctx, &completion.CreateParameters{
			N:         1,
			MaxTokens: 200,
			Prompt:    []string{promptPrefix + text}})
		if err != nil {
			log.Fatalf("Failed to create completion error: %v", err)
		}

		for _, choice := range resp.Choices {
			keywordList = append(keywordList, choice.Text)
		}
	}

	saveKeywords(keywordList)
}

func saveKeywords(keywords []string) {
	f, err := os.Create("keywords.txt")
	if err != nil {
		log.Fatalf("Failed to create keywords file error: %v", err)
	}
	defer f.Close()

	var filteredKeywords []string
	for _, keyword := range keywords {
		if strings.TrimSpace(keyword) != "" {
			filteredKeywords = append(filteredKeywords, keyword)
		}
	}

	// Write filtered array to file
	data := []byte(strings.Join(filteredKeywords, "\n"))
	err = os.WriteFile(f.Name(), data, 0644)
	if err != nil {
		log.Printf("Failed to write keywords to file error: %v", err)
		return
	}

	readKeywordsFile(f)
}

func readKeywordsFile(f *os.File) {
	textData, err := os.ReadFile(f.Name())
	if err != nil {
		log.Printf("Failed to read keywords file error: %v", err)
	}

	// Remove empty lines and "Keywords:" prefix
	lines := strings.Split(string(textData), "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			line = strings.TrimPrefix(line, "Keywords:")
			filteredLines = append(filteredLines, strings.TrimSpace(line))
		}
	}

	fmt.Println(filteredLines)
}
