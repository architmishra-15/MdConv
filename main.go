package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const CurrentVersion = "0.0.2"

const mainTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ .Title }}</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      line-height: 1.6;
      max-width: 800px;
      margin: 0 auto;
      padding: 2rem;
      {{if .BlackBackground}}background-color: #000; color: #fff;{{else}}background-color: #fff; color: #333;{{end}}
    }
    p code, li code, blockquote code {
      background-color: {{if .BlackBackground}}#2d2d2d{{else}}#f1f3f4{{end}};
      color: {{if .BlackBackground}}#f8f8f2{{else}}#d73a49{{end}};
      padding: 0.2em 0.4em;
      border-radius: 3px;
      font-size: 0.9em;
    }
  </style>

  <!-- Highlight.js core (offline) -->
  <script src="./lang_highlight/highlight.min.js"></script>

  <!-- Theme stylesheet -->
  <link rel="stylesheet" href="./lang_highlight/styles/{{.StyleName}}.min.css"
        onerror="this.href='./lang_highlight/styles/github-dark.min.css'">

  <!-- Language definition scripts -->
  {{- range .Langs }}
  <script src="./lang_highlight/languages/{{.}}.min.js"
          onerror="this.src='./lang_highlight/languages/go.min.js'"></script>
  {{- end }}
</head>
<body>
  {{.Content}}
</body>
</html>`

// TemplateData holds data for rendering
// .StyleName       = highlight.js theme name
// .Langs           = list of code block languages to include scripts for
// .Content         = full HTML with injected code blocks
// .BlackBackground = dark mode flag
// .Title           = document title

type TemplateData struct {
	Title           string
	Content         template.HTML
	StyleName       string
	Langs           []string
	BlackBackground bool
}

func main() {
	cfg := parseArgs()
	if err := processMarkdown(cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// Config holds CLI options
type Config struct {
	InputFile       string
	OutputFile      string
	Title           string
	StyleName       string
	BlackBackground bool
}

func parseArgs() *Config {
	cfg := &Config{StyleName: "github-dark", Title: "Markdown Document"}
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: md2html <in.md> [out.html] [--title=] [--style=] [--bg-black]")
		os.Exit(1)
	}
	cfg.InputFile = args[0]
	if len(args) > 1 && !strings.HasPrefix(args[1], "--") {
		cfg.OutputFile = args[1]
	} else {
		base := strings.TrimSuffix(cfg.InputFile, filepath.Ext(cfg.InputFile))
		cfg.OutputFile = base + ".html"
	}
	for _, arg := range args[1:] {
		switch {
		case strings.HasPrefix(arg, "--title="):
			cfg.Title = strings.Trim(strings.TrimPrefix(arg, "--title="), "'\"")
		case strings.HasPrefix(arg, "--style="):
			cfg.StyleName = strings.TrimPrefix(arg, "--style=")
		case arg == "--bg-black":
			cfg.BlackBackground = true
		case arg == "--help" || arg == "-h" || arg == "help":
			Help()
			os.Exit(1)
		case arg == "--version" || arg == "-v" || arg == "version":
			VersionInfo(CurrentVersion)
			os.Exit(1)
		}
	}
	return cfg
}

// processMarkdown reads the markdown, applies GenerateHighlightDiv for each code block,
// and writes the final HTML using mainTemplate
func processMarkdown(cfg *Config) error {
	file, err := os.Open(cfg.InputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	doc, _, err := ParseMarkdown(file)
	if err != nil {
		return err
	}

	// Detect code block languages
	langSet := make(map[string]struct{})
	for _, blk := range doc.Blocks {
		if cb, ok := blk.(*CodeBlock); ok && cb.Lang != "" {
			langSet[cb.Lang] = struct{}{}
		}
	}
	langs := make([]string, 0, len(langSet))
	for l := range langSet {
		langs = append(langs, l)
	}

	var buf bytes.Buffer
	idCounter := 0
	for _, blk := range doc.Blocks {
		switch node := blk.(type) {
		case *CodeBlock:
			idCounter++
			htmlDiv, err := GenerateHighlightDiv(node.Lang, node.Source, cfg.StyleName, fmt.Sprintf("codeblock-%d", idCounter))
			if err != nil {
				node.Render(&buf)
			} else {
				buf.WriteString(string(htmlDiv))
			}
		default:
			blk.Render(&buf)
		}
	}

	tmplData := TemplateData{
		Title:           cfg.Title,
		Content:         template.HTML(buf.String()),
		StyleName:       cfg.StyleName,
		Langs:           langs,
		BlackBackground: cfg.BlackBackground,
	}

	out, err := os.Create(cfg.OutputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	tmpl := template.Must(template.New("main").Parse(mainTemplate))
	return tmpl.Execute(out, tmplData)
}
