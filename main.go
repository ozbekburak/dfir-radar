package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
	"github.com/rakyll/openai-go/completion"
)

var (
	getKeywordsInstruction    = "I will give you an articles about DFIR in various categories Please extract maximum three keywords from the articles to create insights. Use the following format: Keywords: keyword1, keyword2, keyword3. Use only lowercase letters. Do not add dot at the end of the lines.\n"
	createCategoryInstruction = "Please create 4 (four) categories using the keywords below. Decide category names yourself. Just category name and keywords do not add prefix like 'Category includes'. Please create maximum 4 (four) categories. Use the following example as a format: 1. Cybersecurity: NordVPN, LastPass. Use only lowercase letters. Do not use bullet points. Use one line for each category and its keywords.\n"
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

	text, err := analyzeCategories(keywordFile)
	if err != nil {
		log.Printf("Failed to analyze keywords error: %v", err)
		return
	}

	categoryKeywordPair := extractCategories(text)

	csvFile, _, err := saveToCSV(categoryKeywordPair)
	if err != nil {
		log.Printf("Failed to init csv error: %v", err)
		return
	}

	fmt.Println("csvFile name:", csvFile.Name())
}

func extractCategories(text string) map[string][]string {
	// Split the text into lines
	lines := strings.Split(text, "\n")

	// Define a map to store the categories
	categories := make(map[string][]string)

	// Loop over the lines and extract the categories
	for _, line := range lines {
		// Trim leading and trailing whitespace from the line
		line = strings.TrimSpace(line)

		// Check if the line starts with a number followed by a period and a space
		if strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "2.") || strings.HasPrefix(line, "3.") || strings.HasPrefix(line, "4.") {
			// Split the line into the category name and the category items
			parts := strings.SplitN(line, ":", 2)

			// Trim leading and trailing whitespace from the category name and items
			category := strings.TrimSpace(strings.Trim(parts[0], "0123456789."))
			items := strings.TrimSpace(strings.Trim(parts[1], "."))

			// Split the items into individual keywords
			keywords := strings.Split(items, ", ")

			// Add the keywords to the categories map
			categories[category] = keywords
		}
	}

	return categories
}
func saveToCSV(catKeywordPair map[string][]string) (*os.File, *csv.Writer, error) {
	timestamp := time.Now().Format("20060102150405")

	csvFile, err := os.Create(filepath.Join("reports", fmt.Sprintf("report-%s.csv", timestamp)))
	if err != nil {
		log.Printf("Failed to create keywords file error: %v", err)
		return nil, nil, err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write csvHeaders to CSV file
	csvHeaders := make([]string, 0, len(catKeywordPair))
	for category := range catKeywordPair {
		csvHeaders = append(csvHeaders, category)
	}

	if err = csvWriter.Write(csvHeaders); err != nil {
		log.Printf("Failed to write CSV header error: %v", err)
		csvFile.Close()
		return nil, nil, err
	}

	// Write data to CSV file
	numRows := 0
	for _, values := range catKeywordPair {
		if len(values) > numRows {
			numRows = len(values)
		}
	}
	for i := 0; i < numRows; i++ {
		record := make([]string, 0, len(catKeywordPair))
		for _, category := range csvHeaders {
			if len(catKeywordPair[category]) <= i {
				record = append(record, "") // add a placeholder value for missing keywords
			} else {
				record = append(record, catKeywordPair[category][i])
			}
		}
		if err = csvWriter.Write(record); err != nil {
			log.Printf("Failed to write CSV record error: %v", err)
		}
	}

	return csvFile, csvWriter, nil
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
		answers, err := askChatGPT(getKeywordsInstruction + text)
		if err != nil {
			log.Printf("Failed to ask chat gpt error: %v", err)
			return nil, err
		}

		keywordList = append(keywordList, answers...)
	}

	return keywordList, nil
}

func askChatGPT(prompt string) ([]string, error) {
	client := chat.NewClient(openai.NewSession(os.Getenv("OPENAI_API_KEY")), "gpt-3.5-turbo")
	resp, err := client.CreateCompletion(context.Background(), &chat.CreateCompletionParams{
		Messages: []*chat.Message{
			{
				Role:    "user",
				Content: prompt},
		},
	})
	if err != nil {
		log.Printf("Failed to create openai gpt-3.5-turbo completion error: %v", err)
		// If we exceeded the limit, try to ask davinci model
		if strings.Contains(err.Error(), "status_code=429") {
			log.Println("Exceeded the limit, trying to ask davinci model...")
			answers, err := askDavinci([]string{prompt})
			if err != nil {
				log.Printf("Failed to ask chat gpt with model davinci-text-003 error: %v", err)
				return nil, err
			}
			var keywordList []string
			for _, answer := range answers {
				keywordList = append(keywordList, answer.Text)
			}

			return keywordList, nil
		}
		return nil, err
	}

	var contents []string
	for _, choice := range resp.Choices {
		contents = append(contents, choice.Message.Content)
	}

	return contents, nil
}

func askDavinci(prompt []string) ([]*completion.Choice, error) {
	client := completion.NewClient(openai.NewSession(os.Getenv("OPENAI_API_KEY")), "text-davinci-003")
	resp, err := client.Create(context.Background(), &completion.CreateParams{
		N:         1,
		MaxTokens: 200,
		Prompt:    prompt})
	if err != nil {
		log.Printf("Failed to create openai completion error with model of davinci-text-003: %v", err)
		return nil, err
	}

	return resp.Choices, nil
}

func saveKeywords(keywords []string) (*os.File, error) {
	timestamp := time.Now().Format("20060102150405")

	f, err := os.Create(filepath.Join("keywords", fmt.Sprintf("keywords-%s.txt", timestamp)))
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

func analyzeCategories(f *os.File) (string, error) {
	keywordList, err := readKeywordsFile(f)
	if err != nil {
		log.Printf("Failed to read keywords file error: %v", err)
		return "", err
	}

	categories, err := askChatGPT(createCategoryInstruction + strings.Join(keywordList, "\n"))
	if err != nil {
		log.Printf("Failed to ask chat gpt categories error: %v", err)
		return "", err
	}

	var category string
	for _, category = range categories {
		log.Println(category)
	}

	return category, nil
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
