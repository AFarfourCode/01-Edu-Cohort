package ascii

import (
	"os"
	"strings"
)

// GenerateArt returns the ASCII art as a string instead of printing it.
// دالة تقوم بتوليد الرسم وإعادته كنص بدلاً من طباعته
func GenerateArt(input string, fontFile string) (string, error) {
	// 1. Read the banner file
	// قراءة ملف الخط
	path := "banners/" + fontFile + ".txt"
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// 2. Normalize line endings (Windows fix)
	fileContent := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(fileContent, "\n")

	// 3. Prepare input
	input = strings.ReplaceAll(input, "\r\n", "\n")
	inputLines := strings.Split(input, "\n")

	var resultBuilder strings.Builder

	// 4. Generate Logic
	for _, word := range inputLines {
		if word == "" {
			resultBuilder.WriteString("\n")
			continue
		}

		for i := 0; i < 8; i++ {
			for _, ch := range word {
				startLine := int(ch-32)*9 + 1 + i
				if startLine >= 0 && startLine < len(lines) {
					resultBuilder.WriteString(lines[startLine])
				}
			}
			resultBuilder.WriteString("\n")
		}
	}

	return resultBuilder.String(), nil
}
