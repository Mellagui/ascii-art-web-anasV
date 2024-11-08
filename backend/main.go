package main

import (
	"ascii-art-web/utils"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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
}

func init() {
    styles := []string{"standard", "shadow", "thinkertoy", "hy"}
    for _, style := range styles {
        bannerFile := filepath.Join("banners", style+".txt")
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
    http.HandleFunc("/generate", handleGenerate)
    http.HandleFunc("/download", handleExport)

    log.Println("Server starting on : http://localhost:8080/")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
        http.Error(w, "404- Not Found", 404)
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
	default :
		http.Error(w, "404- Not Found", 404)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    renderTemplate(w, PageData{Title: "ASCII Art Generator"})
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    inputText := r.FormValue("text")
    style := r.FormValue("style")

    if strings.TrimSpace(inputText) == "" {
        renderTemplate(w, PageData{Title: "ASCII Art Generator", Error: "Text cannot be empty"})
        return
    }

    bannerMap, exists := bannerMaps[style]
    if !exists {
        renderTemplate(w, PageData{Title: "ASCII Art Generator", Error: "Invalid style selected"})
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
    
    tmplPath := filepath.Join("../", "templates", "index.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
