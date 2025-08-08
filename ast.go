// ast.go

package main

import (
	"fmt"
	"html"
	"io"
	"strings"
)

type Node interface {
	Render(w io.Writer) error
}

type Document struct {
	Blocks []Node
	CSS    string
}

type Heading struct {
	Level int
	Text  string
}

type Paragraph struct {
	Lines []string
}

// Blockquote for quoted text (> text)
type Blockquote struct {
	Lines []string
}

// unordered list
type List struct {
	Items []string
}

type CodeBlock struct {
	Lang            string
	Source          []string
	HighlightedHTML string // Highlighted HTML obtained from Chroma
	IsHighlighted   bool   // Flag to check if highlighting was applied
}

func (d *Document) Render(w io.Writer) error {
	for _, blk := range d.Blocks {
		if err := blk.Render(w); err != nil {
			return err
		}
	}
	return nil
}

func (d *Document) AddCSS(css string) {
	if d.CSS == "" {
		d.CSS = css
	}
	if !strings.Contains(d.CSS, css) {
		d.CSS += "\n" + css
	}
}

func (h *Heading) Render(w io.Writer) error {
	// Text already contains processed HTML from inline formatting, so write directly
	_, err := fmt.Fprintf(w, "<h%d>%s</h%d>\n", h.Level, h.Text, h.Level)
	return err
}

func (p *Paragraph) Render(w io.Writer) error {
	// Lines already contain processed HTML from inline formatting
	// Join them without additional escaping and write directly
	content := strings.Join(p.Lines, " ")
	_, err := fmt.Fprintf(w, "<p>%s</p>\n", content)
	return err
}

func (b *Blockquote) Render(w io.Writer) error {
	io.WriteString(w, "<blockquote>\n")

	// Group lines into paragraphs (split by empty lines)
	paragraphs := [][]string{{}}
	currentParagraph := 0

	for _, line := range b.Lines {
		if line == "" {
			if len(paragraphs[currentParagraph]) > 0 {
				paragraphs = append(paragraphs, []string{})
				currentParagraph++
			}
		} else {
			paragraphs[currentParagraph] = append(paragraphs[currentParagraph], line)
		}
	}

	// Render each paragraph - content already has HTML formatting
	for _, paragraph := range paragraphs {
		if len(paragraph) > 0 {
			content := strings.Join(paragraph, " ")
			fmt.Fprintf(w, "  <p>%s</p>\n", content)
		}
	}

	io.WriteString(w, "</blockquote>\n")
	return nil
}

func (l *List) Render(w io.Writer) error {
	io.WriteString(w, "<ul>\n")

	for _, item := range l.Items {
		// Items already contain processed HTML from inline formatting
		fmt.Fprintf(w, "  <li>%s</li>\n", item)
	}

	io.WriteString(w, "</ul>\n")
	return nil
}

func (c *CodeBlock) Render(w io.Writer) error {
	// Wrap the code block in a container with code-block class
	io.WriteString(w, `<div class="code-block">`)

	// render highlighted HTML directly
	if c.IsHighlighted && c.HighlightedHTML != "" {
		_, err := io.WriteString(w, c.HighlightedHTML)
		if err != nil {
			return err
		}
	} else {
		// Fallback if no highlighting available
		io.WriteString(w, "<pre><code")
		if c.Lang != "" {
			io.WriteString(w, ` class="language-`+c.Lang+`"`)
		}

		io.WriteString(w, ">")
		io.WriteString(w, html.EscapeString(strings.Join(c.Source, "\n")))
		io.WriteString(w, "</code></pre>")
	}

	// Close the wrapper div
	io.WriteString(w, "</div>\n")
	return nil
}

func (c *CodeBlock) SetHighlightedContent(HighlightedHTML string) {
	c.HighlightedHTML = HighlightedHTML
	c.IsHighlighted = true
}
