package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/completion"
)

var (
	getKeywordsInstruction    = "I will give you an articles about DFIR in various categories Please extract maximum three keywords from the articles to create insights. Use the following format: Keywords: keyword1, keyword2, keyword3. Use only lowercase letters \n"
	createCategoryInstruction = "Please create 4 categories using the keywords below. Decide category names yourself. Please create maximum 4 (four) categories.\n"
	newsDir                   = "news"
	textList                  []string
)

func main() {
	texts, err := readFiles(newsDir)
	if err != nil {
		log.Printf("Failed to read files from directory '%s' error: %v", newsDir, err)
		return
	}

	keywords, err := keywords(texts)
	if err != nil {
		log.Printf("Failed to get keywords from chatGPT error: %v", err)
		return
	}

	keywordFile, err := saveKeywords(keywords)
	if err != nil {
		log.Printf("Failed to save keywords error: %v", err)
		return
	}

	if err = analyzeCategories(keywordFile); err != nil {
		log.Printf("Failed to analyze keywords error: %v", err)
		return
	}
}

func readFiles(dir string) ([]string, error) {
	newsFiles, err := os.ReadDir(newsDir)
	if err != nil {
		log.Printf("Failed to read news directory '%s' error: %v", newsDir, err)
		return nil, err
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

	return textList, nil
}

func keywords(texts []string) ([]string, error) {

	var keywordList []string
	for _, text := range texts {
		answers, err := askChatGPT([]string{getKeywordsInstruction + text})
		if err != nil {
			log.Printf("Failed to ask chat gpt error: %v", err)
			return nil, err
		}

		for _, answer := range answers {
			keywordList = append(keywordList, answer.Text)
		}
	}

	return keywordList, nil
}

func askChatGPT(prompt []string) ([]*completion.Choice, error) {
	client := completion.NewClient(openai.NewSession(os.Getenv("OPENAI_API_KEY")), "text-davinci-003")
	resp, err := client.Create(context.Background(), &completion.CreateParameters{
		N:         1,
		MaxTokens: 200,
		Prompt:    prompt})
	if err != nil {
		log.Printf("Failed to create openai completion error: %v", err)
		return nil, err
	}

	return resp.Choices, nil
}

func saveKeywords(keywords []string) (*os.File, error) {
	timestamp := time.Now().Format("20060102150405")

	f, err := os.Create(fmt.Sprintf("keywords-%s.txt", timestamp))
	if err != nil {
		log.Printf("Failed to create keywords file error: %v", err)
		return nil, err
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
	if err = os.WriteFile(f.Name(), data, 0644); err != nil {
		log.Printf("Failed to write keywords to file error: %v", err)
		return nil, err
	}

	return f, nil
}

func analyzeCategories(f *os.File) error {
	keywordList, err := readKeywordsFile(f)
	if err != nil {
		log.Printf("Failed to read keywords file error: %v", err)
		return err
	}

	categories, err := askChatGPT([]string{createCategoryInstruction + strings.Join(keywordList, "\n")})
	if err != nil {
		log.Printf("Failed to ask chat gpt categories error: %v", err)
		return err
	}

	for i, category := range categories {
		fmt.Printf("category %d: %s\n", i, category.Text)
	}
	return nil
}

func readKeywordsFile(f *os.File) ([]string, error) {
	textData, err := os.ReadFile(f.Name())
	if err != nil {
		log.Printf("Failed to read keywords file error: %v", err)
		return nil, err
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

	return filteredLines, nil
}
