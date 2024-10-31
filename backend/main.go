package main

import (
	"ascii-art-web/utils"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
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
	fs := http.FileServer(http.Dir("../static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Route handlers
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/generate", handleGenerate)

    log.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
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