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

const mainTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <style>
    /* Base styles for better typography and layout */
    body { 
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      line-height: 1.6;
      max-width: 800px;
      margin: 0 auto;
      padding: 2rem;
      color: #333;
    }
    pre {
      overflow-x: auto;
      padding: 1rem;
      border-radius: 8px;
      margin: 1rem 0;
    }
    code {
      font-family: "Fira Code", "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, monospace;
    }
    
    /* Inline code styling */
    p code, li code, h1 code, h2 code, h3 code, h4 code, h5 code, h6 code, blockquote code {
      background-color: #f1f3f4;
      padding: 0.2em 0.4em;
      border-radius: 3px;
      font-size: 0.9em;
      color: #d73a49;
      border: 1px solid #e1e4e8;
    }
    
    /* Blockquote styling */
    blockquote {
      margin: 1rem 0;
      padding: 0 1rem;
      border-left: 4px solid #dfe2e5;
      background-color: #f8f9fa;
      color: #6a737d;
    }
    blockquote p {
      margin: 0.5rem 0;
    }
    blockquote p:first-child {
      margin-top: 0;
    }
    blockquote p:last-child {
      margin-bottom: 0;
    }
    
    h1, h2, h3, h4, h5, h6 {
      margin-top: 2rem;
      margin-bottom: 1rem;
    }
    p {
      margin-bottom: 1rem;
    }
    ul {
      margin-bottom: 1rem;
    }
    
    /* Syntax highlighting styles will be injected here */
    {{.CSS}}
  </style>
  <title>{{.Title}}</title>
</head>
<body {{if .BodyClass}}class="{{.BodyClass}}"{{end}}>
  {{.Content}}
</body>
</html>
`

// TemplateData holds data for the HTML template
type TemplateData struct {
	Title     string
	CSS       template.CSS
	Content   template.HTML
	BodyClass string
}

// Config holds configuration options
type Config struct {
	InputFile  string
	OutputFile string
	Title      string
	StyleName  string
	BodyClass  string
}

func main() {
	config := parseArgs()

	if err := processMarkdown(config); err != nil {
		log.Fatalf("Error processing markdown: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", config.InputFile, config.OutputFile)
}

// parseArgs parses command line arguments (simplified version)
func parseArgs() *Config {
	config := &Config{
		StyleName: "github", // default syntax highlighting style
		Title:     "Markdown Document",
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: go run . <input.md> [output.html] [--style=github] [--title='My Doc']")
		fmt.Println("Available styles:", strings.Join(GetAvailableStyles(), ", "))
		os.Exit(1)
	}

	config.InputFile = args[0]

	// Set default output file
	if len(args) > 1 && !strings.HasPrefix(args[1], "--") {
		config.OutputFile = args[1]
	} else {
		base := strings.TrimSuffix(config.InputFile, filepath.Ext(config.InputFile))
		config.OutputFile = base + ".html"
	}

	// Parse additional flags
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--style=") {
			config.StyleName = strings.TrimPrefix(arg, "--style=")
		} else if strings.HasPrefix(arg, "--title=") {
			config.Title = strings.Trim(strings.TrimPrefix(arg, "--title="), "\"'")
		} else if strings.HasPrefix(arg, "--body-class=") {
			config.BodyClass = strings.TrimPrefix(arg, "--body-class=")
		}
	}

	return config
}

// processMarkdown handles the entire markdown to HTML conversion process
func processMarkdown(config *Config) error {
	// Read input file
	file, err := os.Open(config.InputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer file.Close()

	// Parse markdown into AST
	doc, err := ParseMarkdown(file)
	if err != nil {
		return fmt.Errorf("error parsing markdown: %w", err)
	}

	// Apply syntax highlighting
	highlighter := NewSyntaxHighlighter(config.StyleName)
	if err := highlighter.ProcessCodeBlocks(doc); err != nil {
		return fmt.Errorf("error processing code blocks: %w", err)
	}

	// Render document to HTML
	var contentBuf bytes.Buffer
	if err := doc.Render(&contentBuf); err != nil {
		return fmt.Errorf("error rendering document: %w", err)
	}

	// Prepare template data
	templateData := TemplateData{
		Title:     config.Title,
		CSS:       template.CSS(doc.CSS),
		Content:   template.HTML(contentBuf.String()),
		BodyClass: config.BodyClass,
	}

	// Create output file
	outputFile, err := os.Create(config.OutputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	// Parse and execute template
	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	if err := tmpl.Execute(outputFile, templateData); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
