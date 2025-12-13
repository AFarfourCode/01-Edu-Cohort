package main

import (
	"html/template" // Used to parse and execute HTML files (render dynamic data)
	"log"           // Used for printing errors and server status to the terminal
	"net/http"      // The core package for creating the web server and handling requests
	"os"            // Used for file system operations (checking if files exist, reading files)
	"strings"       // Used for string manipulation (splitting, replacing)
)

// --- STRUCTS (Data Containers) ---

// PageData holds the information we want to send to the HTML file (index.html).
// "Input" is the text the user typed, "Result" is the ASCII art we generated.
type PageData struct {
	Result string
	Input  string
}

// ErrorData holds information for the error page (error.html).
// It allows us to show dynamic error messages (like 404 or 500).
type ErrorData struct {
	StatusCode int
	Message    string
}

// --- MAIN FUNCTION (Server Entry Point) ---

func main() {
	// 1. Serve Static Files (CSS/Images)
	// We tell Go: "If a URL starts with /static/, look inside the ./static folder."
	// http.StripPrefix removes "/static/" from the URL so we look for "style.css" instead of "static/style.css" inside the folder.
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Register Routes (Endpoints)
	// When someone visits the home page "/", run the 'homeHandler' function.
	http.HandleFunc("/", homeHandler)
	// When someone submits the form to "/ascii-art", run the 'asciiHandler' function.
	http.HandleFunc("/ascii-art", asciiHandler)

	// 3. Start the Server
	// We log a message so you know the server is running.
	log.Println("Server starting on http://localhost:8080")

	// ListenAndServe starts the web server on port 8080.
	// It blocks the program here (keeps it running) until you stop it (Ctrl+C).
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err) // If the server crashes (e.g., port busy), print error and exit.
	}
}

// --- HANDLERS (Controllers) ---

// homeHandler manages the main page request.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Security/Logic Check: Ensure the URL is exactly "/".
	// Without this, "/" acts as a catch-all for any unknown page not defined elsewhere.
	if r.URL.Path != "/" {
		renderError(w, http.StatusNotFound, "PAGE NOT FOUND")
		return
	}

	// Method Check: We only want GET requests here (viewing the page).
	if r.Method != http.MethodGet {
		renderError(w, http.StatusMethodNotAllowed, "METHOD NOT ALLOWED")
		return
	}

	// If everything is okay, render the index.html file.
	// We pass 'nil' because we don't have any data to show yet (first visit).
	renderTemplate(w, "index.html", nil)
}

// asciiHandler manages the form submission.
func asciiHandler(w http.ResponseWriter, r *http.Request) {
	// Method Check: Forms submit data using POST. If it's not POST, reject it.
	if r.Method != http.MethodPost {
		renderError(w, http.StatusMethodNotAllowed, "METHOD NOT ALLOWED")
		return
	}

	// ParseForm reads the data sent by the browser.
	if err := r.ParseForm(); err != nil {
		renderError(w, http.StatusBadRequest, "BAD REQUEST")
		return
	}

	// Retrieve the values from the HTML inputs using their 'name' attributes.
	text := r.FormValue("text")     // <input name="text">
	banner := r.FormValue("banner") // <select name="banner">

	// Validation: Ensure text is not empty.
	if text == "" {
		renderError(w, http.StatusBadRequest, "MISSING INPUT TEXT")
		return
	}

	// Validation: Ensure the banner is one of the allowed types.
	// This prevents hackers from trying to read other files on your computer.
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		renderError(w, http.StatusBadRequest, "INVALID BANNER STYLE")
		return
	}

	// Call our logic function to create the art.
	asciiArt, err := GenerateAsciiArt(text, banner)
	if err != nil {
		log.Printf("Error generating art: %v", err) // Log error to console for the developer
		renderError(w, http.StatusInternalServerError, "INTERNAL SERVER ERROR")
		return
	}

	// Prepare the data to send back to the HTML.
	// We send 'Result' (the art) and 'Input' (so the text box doesn't go blank).
	data := PageData{
		Result: asciiArt,
		Input:  text,
	}
	renderTemplate(w, "index.html", data)
}

// --- HELPERS (Utilities) ---

// renderTemplate handles the repetitive task of parsing and executing HTML files.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// 1. Check if the file actually exists before trying to parse it.
	if _, err := os.Stat("templates/" + tmpl); os.IsNotExist(err) {
		renderError(w, http.StatusNotFound, "TEMPLATE NOT FOUND")
		return
	}

	// 2. Parse the HTML file.
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "TEMPLATE PARSING ERROR")
		return
	}

	// 3. Execute the template (inject the 'data' into the HTML) and write to 'w'.
	if err := t.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// renderError is a dedicated function to show pretty error pages instead of plain text.
func renderError(w http.ResponseWriter, status int, message string) {
	// Set the HTTP status code (e.g., 404, 500) in the header.
	w.WriteHeader(status)

	// Try to load the custom error.html page.
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// Fallback: If error.html is missing, just send plain text so the user sees something.
		http.Error(w, message, status)
		return
	}

	// Send the status code and message to the HTML to be displayed.
	t.Execute(w, ErrorData{StatusCode: status, Message: message})
}

// --- LOGIC (ASCII Generation) ---

// GenerateAsciiArt contains the core algorithm to convert text to ASCII art.
func GenerateAsciiArt(input string, banner string) (string, error) {
	// Construct the file path (e.g., "banners/standard.txt")
	path := "banners/" + banner + ".txt"

	// Read the entire font file into memory.
	readinput, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Fix line endings: Windows uses \r\n, logic expects \n.
	editInput := strings.ReplaceAll(string(readinput), "\r\n", "\n")
	// Split the file into a slice of strings (each line of the file is an item).
	lines := strings.Split(editInput, "\n")

	// Normalize the user's input text (handle Windows newlines).
	input = strings.ReplaceAll(input, "\r\n", "\n")
	argslines := strings.Split(input, "\n")

	var output strings.Builder

	// Loop through every line of the user's input (in case they typed multiple lines).
	for _, word := range argslines {
		// If the line is empty (user pressed Enter twice), just print a newline.
		if word == "" {
			output.WriteString("\n")
			continue
		}

		// The ASCII art is 8 lines high. We need to loop 8 times to build one row of text.
		for i := 0; i < 8; i++ {
			// Loop through every character in the word (e.g., 'H', 'E', 'L', 'L', 'O').
			for _, ch := range word {
				// MATH EXPLANATION:
				// 'ch' is the ASCII value (e.g., 'A' is 65).
				// We subtract 32 because the font file starts at space (ASCII 32).
				// We multiply by 9 because every character takes up 9 lines in the file.
				// We add 1 to skip the very first empty line in the file.
				// We add 'i' to get the current slice (0 to 7) of the character.
				startline := int(ch-32)*9 + 1 + i

				// Safety check: ensure we don't try to read a line outside the file.
				if startline >= 0 && startline < len(lines) {
					output.WriteString(lines[startline])
				}
			}
			// After printing the slice of all letters, move to the next line.
			output.WriteString("\n")
		}
	}
	return output.String(), nil
}
