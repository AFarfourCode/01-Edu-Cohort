package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageData struct {
	Input  string
	Result string
}

func main() {

	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/ascii-art", asciiHandler)

	log.Println("Server started at http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error (Template missing)", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" || banner == "" {
		http.Error(w, "400 - Bad Request (Missing input)", http.StatusBadRequest)
		return
	}

	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		http.Error(w, "400 - Bad Request (Invalid banner)", http.StatusBadRequest)
		return
	}

	art, err := GenerateAsciiArt(text, banner)
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Input:  text,
		Result: art,
	}

	tmpl.Execute(w, data)
}

func GenerateAsciiArt(input string, banner string) (string, error) {
	path := "banners/" + banner + ".txt"

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	fileContent := strings.ReplaceAll(string(data), "\r\n", "\n")

	lines := strings.Split(fileContent, "\n")

	input = strings.ReplaceAll(input, "\\n", "\n")

	inputLines := strings.Split(input, "\n")

	var output strings.Builder

	for _, line := range inputLines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		for i := 0; i < 8; i++ {
			for _, char := range line {
				start := int(char-32)*9 + 1 + i

				if start >= 0 && start < len(lines) {
					output.WriteString(lines[start])
				}
			}
			output.WriteString("\n")
		}
	}
	return output.String(), nil
}
