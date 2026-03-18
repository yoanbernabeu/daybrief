package newsletter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yoanbernabeu/daybrief/internal/gemini"
)

func SaveJSON(newsletter *gemini.Newsletter, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("creating output dir: %w", err)
	}

	filename := fmt.Sprintf("%s.json", time.Now().Format("2006-01-02"))
	path := filepath.Join(outputDir, filename)

	data, err := json.MarshalIndent(newsletter, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling newsletter: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("writing newsletter: %w", err)
	}

	return path, nil
}

func LoadJSON(path string) (*gemini.Newsletter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading newsletter: %w", err)
	}

	var nl gemini.Newsletter
	if err := json.Unmarshal(data, &nl); err != nil {
		return nil, fmt.Errorf("parsing newsletter: %w", err)
	}

	return &nl, nil
}

func GetLatestOutputPath(outputDir string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(outputDir, "*.json"))
	if err != nil {
		return "", fmt.Errorf("globbing output dir: %w", err)
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no newsletter outputs found in %s", outputDir)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(matches)))
	return matches[0], nil
}
