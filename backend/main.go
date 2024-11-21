package main

import (
	"ascii-art-web/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var bannerMaps = make(map[string]map[rune][]string)

type PageData struct {
	Title     string
	AsciiArt  string
	InputText string
	StyleUsed string
	Error     string
	Status    int
}

func init() {
	styles := []string{"standard", "shadow", "thinkertoy", "hy"}
	for _, style := range styles {
		bannerFile := "banners/" + style + ".txt"
		banner, err := utils.LoadBanner(bannerFile)
		if err != nil {
			log.Fatalf("Failed to load banner %s: %v", style, err)
		}
		bannerMaps[style] = banner
	}
}

func main() {
	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static", http.HandlerFunc(handleStatic)))

	// Route handlers
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ascii-art", handleGenerate)
	http.HandleFunc("/download", handleExport)

	log.Println("Server starting on : http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {

	HandleStylePath := strings.HasPrefix(r.URL.Path, "/styles.css/") //
	if r.URL.Path == "/" || HandleStylePath {
		showError(w, "404- Not Found", 404)
		return
	}
	fs := http.FileServer(http.Dir("../static"))
	fs.ServeHTTP(w, r)
}

func handleExport(w http.ResponseWriter, r *http.Request) {
	fileType := r.FormValue("fileType")
	asciiArt := r.FormValue("asciiArt")
	switch fileType {
	case "txt":
		w.Header().Set("Content-Disposition", "attachment; filename=asciiArt.txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Lenght", strconv.Itoa(len(asciiArt)))
		w.Write([]byte(asciiArt))
	case "doc":
		w.Header().Set("Content-Disposition", "attachment; filename=asciiArt.docx")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Lenght", strconv.Itoa(len(asciiArt)))
		w.Write([]byte(asciiArt))
	case "html":
		result := "<pre>" + asciiArt + "</pre>"
		w.Header().Set("Content-Disposition", "attachment; filename=asciiArt.html")
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Lenght", strconv.Itoa(len(result)))
		w.Write([]byte(result))
	default:
		showError(w, "404- Not Found", 404)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		showError(w, "404 - not Found", 404)
		return
	}
	renderTemplate(w, PageData{Title: "ASCII Art Generator"})
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		showError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inputText := r.FormValue("text")
	style := r.FormValue("style")

	if strings.TrimSpace(inputText) == "" {
		renderTemplate(w, PageData{Title: "ASCII Art Generator", Error: "Text cannot be empty", Status: http.StatusBadRequest})
		return
	}

	if len(inputText) > 2000 {
		renderTemplate(w, PageData{Title: "ASCII Art Generator", Error: "Input text must be less than 2000 charcters", Status: http.StatusBadRequest})
		return
	}

	bannerMap, exists := bannerMaps[style]
	if !exists {
		renderTemplate(w, PageData{Title: "ASCII Art Generator", Error: "Invalid style selected", Status: http.StatusBadRequest})
		return
	}

	var sb strings.Builder
	utils.PrintAsciiArt(inputText, bannerMap, "", &sb)
	renderTemplate(w, PageData{
		Title:     "ASCII Art Generator",
		AsciiArt:  sb.String(),
		InputText: inputText,
		StyleUsed: style,
	})
}

func renderTemplate(w http.ResponseWriter, data PageData) {

	//tmplPath := filepath.Join("../", "templates", "index.html")
	if data.Status != 0 {
		w.WriteHeader(data.Status)
	}
	tmpl, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		showError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		showError(w, err.Error(), http.StatusInternalServerError)
	}
}

// ----------- HTML ERROR -----------

type Error struct {
	Status  int
	Message string
}

// Function to render error pages with an HTTP status code
func showError(w http.ResponseWriter, message string, status int) {

	// Set the HTTP status code
	w.WriteHeader(status)

	// Parse the error template
	tmpl, err := template.ParseFiles("../templates/ErrPage.html")
	if err != nil {
		// If template parsing fails, fallback to a generic error response
		http.Error(w, "Could not load error page", http.StatusInternalServerError)
		return
	}

	httpError := Error{
		Status:  status,
		Message: message,
	}
	// Execute the template with the error message
	err = tmpl.Execute(w, httpError)
	if err != nil {
		// If template execution fails, respond with a generic error
		http.Error(w, "Could not render error page", http.StatusInternalServerError)
	}
}
