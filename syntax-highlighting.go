package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var (
	// Regex to extract CSS from <style> tags in Chroma output
	cssRegex = regexp.MustCompile(`<style[^>]*>(.*?)</style>`)
	// Regex to extract body content (everything between <pre> tags or similar)
	bodyRegex = regexp.MustCompile(`<pre[^>]*>.*?</pre>`)
)

// SyntaxHighlighter handles all syntax highlighting operations
type SyntaxHighlighter struct {
	styleName string
	formatter chroma.Formatter
	style     *chroma.Style
	cssCache  string // Cache CSS to avoid regenerating
}

// NewSyntaxHighlighter creates a new syntax highlighter with given style
func NewSyntaxHighlighter(styleName string) *SyntaxHighlighter {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Get("github") // fallback to github style
	}

	// Create HTML formatter with CSS classes
	formatter := formatters.Get("html")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	return &SyntaxHighlighter{
		styleName: styleName,
		formatter: formatter,
		style:     style,
	}
}

// getCSS generates CSS for the current style (cached)
func (sh *SyntaxHighlighter) getCSS() string {
	if sh.cssCache != "" {
		return sh.cssCache
	}

	// Generate CSS by creating a full HTML output with a simple example
	// and extracting the CSS from it
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, `package main`)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer

	// Try to create an HTML formatter that includes CSS
	htmlFormatter := formatters.Get("html")
	if htmlFormatter != nil {
		// Format with the iterator
		err = htmlFormatter.Format(&buf, sh.style, iterator)
		if err != nil {
			return ""
		}

		// Extract CSS from the output
		output := buf.String()
		if cssMatches := cssRegex.FindStringSubmatch(output); len(cssMatches) > 1 {
			sh.cssCache = strings.TrimSpace(cssMatches[1])
			return sh.cssCache
		}
	}

	// Fallback: generate basic CSS manually for common token types
	sh.cssCache = sh.generateBasicCSS()
	return sh.cssCache
}

// generateBasicCSS creates basic CSS for syntax highlighting
func (sh *SyntaxHighlighter) generateBasicCSS() string {
	// This is a fallback CSS based on common token types
	// You can customize this based on your needs
	return `
/* Syntax highlighting */
.chroma { background-color: #f8f8f8; }
.chroma .err { color: #a61717; background-color: #e3d2d2; }
.chroma .k { color: #000000; font-weight: bold; }
.chroma .kd { color: #000000; font-weight: bold; }
.chroma .kn { color: #000000; font-weight: bold; }
.chroma .kp { color: #000000; font-weight: bold; }
.chroma .kr { color: #000000; font-weight: bold; }
.chroma .kt { color: #445588; font-weight: bold; }
.chroma .nc { color: #445588; font-weight: bold; }
.chroma .no { color: #008080; }
.chroma .nd { color: #3c5d5d; font-weight: bold; }
.chroma .ni { color: #800080; }
.chroma .ne { color: #990000; font-weight: bold; }
.chroma .nf { color: #990000; font-weight: bold; }
.chroma .nl { color: #990000; font-weight: bold; }
.chroma .nn { color: #555555; }
.chroma .nt { color: #000080; }
.chroma .nv { color: #008080; }
.chroma .s { color: #dd1144; }
.chroma .sa { color: #dd1144; }
.chroma .sb { color: #dd1144; }
.chroma .sc { color: #dd1144; }
.chroma .sd { color: #dd1144; }
.chroma .s2 { color: #dd1144; }
.chroma .se { color: #dd1144; }
.chroma .sh { color: #dd1144; }
.chroma .si { color: #dd1144; }
.chroma .sx { color: #dd1144; }
.chroma .sr { color: #009926; }
.chroma .s1 { color: #dd1144; }
.chroma .ss { color: #990073; }
.chroma .m { color: #009999; }
.chroma .mf { color: #009999; }
.chroma .mh { color: #009999; }
.chroma .mi { color: #009999; }
.chroma .mo { color: #009999; }
.chroma .c { color: #999988; font-style: italic; }
.chroma .ch { color: #999988; font-style: italic; }
.chroma .cm { color: #999988; font-style: italic; }
.chroma .c1 { color: #999988; font-style: italic; }
.chroma .cs { color: #999999; font-weight: bold; font-style: italic; }
`
}

// HighlightCode highlights a single code block and returns CSS and HTML separately
func (sh *SyntaxHighlighter) HighlightCode(code, language string) (css string, html string, err error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", "", err
	}

	// Get the highlighted HTML
	var htmlBuf bytes.Buffer
	err = sh.formatter.Format(&htmlBuf, sh.style, iterator)
	if err != nil {
		return "", "", err
	}

	// Get CSS (cached)
	css = sh.getCSS()

	return css, htmlBuf.String(), nil
}

// ExtractCSSAndBody extracts CSS and body content from Chroma's full HTML output
// This is an alternative method if you get full HTML from Chroma
func ExtractCSSAndBody(chromaHTML string) (css string, body string) {
	// Extract CSS
	cssMatches := cssRegex.FindStringSubmatch(chromaHTML)
	if len(cssMatches) > 1 {
		css = strings.TrimSpace(cssMatches[1])
	}

	// Extract body content (the actual highlighted code)
	bodyMatches := bodyRegex.FindAllString(chromaHTML, -1)
	if len(bodyMatches) > 0 {
		body = strings.Join(bodyMatches, "\n")
	}

	return css, body
}

// ProcessCodeBlocks processes all code blocks in a document and adds syntax highlighting
func (sh *SyntaxHighlighter) ProcessCodeBlocks(doc *Document) error {
	var hasCodeBlocks bool

	for _, block := range doc.Blocks {
		if codeBlock, ok := block.(*CodeBlock); ok {
			if codeBlock.Lang != "" && len(codeBlock.Source) > 0 {
				code := strings.Join(codeBlock.Source, "\n")

				_, html, err := sh.HighlightCode(code, codeBlock.Lang)
				if err != nil {
					// If highlighting fails, leave the code block as is
					fmt.Printf("Warning: Failed to highlight %s code block: %v\n", codeBlock.Lang, err)
					continue
				}

				// Set the highlighted content
				codeBlock.SetHighlightedContent(html)
				hasCodeBlocks = true
			}
		}
	}

	// Add CSS only once if we have code blocks
	if hasCodeBlocks {
		css := sh.getCSS()
		if css != "" {
			doc.AddCSS(css)
		}
	}

	return nil
}

// GetAvailableStyles returns a list of available Chroma styles
func GetAvailableStyles() []string {
	var styleNames []string
	for _, style := range styles.Registry {
		styleNames = append(styleNames, style.Name)
	}
	return styleNames
}

// GetAvailableLanguages returns a list of available lexers/languages
func GetAvailableLanguages() []string {
	var languages []string
	for _, name := range lexers.Names(true) {
		lexer := lexers.Get(name)
		if lexer != nil {
			languages = append(languages, lexer.Config().Name)
		}
	}
	return languages
}
