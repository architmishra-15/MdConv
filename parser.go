package main

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var (
	headingRegex    = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	codeBlockRegex  = regexp.MustCompile("^```(\\w*)\\s*$")
	listItemRegex   = regexp.MustCompile(`^\s*[-*+]\s+(.+)$`)
	blockquoteRegex = regexp.MustCompile(`^\s*>\s*(.*)$`)
	inlineCodeRegex = regexp.MustCompile("`([^`]+)`")

	boldRegex        = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	boldItalicRegex  = regexp.MustCompile(`\*\*\*([^*]+)\*\*\*`)
	italicStarRegex  = regexp.MustCompile(`\*([^*]+)\*`)
	italicUnderRegex = regexp.MustCompile(`_([^_]+)_`)
)

type Parser struct {
	scanner *bufio.Scanner
	current string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
	}
}

func (p *Parser) nextLine() bool {
	if p.scanner.Scan() {
		p.current = p.scanner.Text()
		return true
	}

	p.current = ""
	return false
}

func (p *Parser) skipEmptyLines() {
	for p.current == "" {
		p.nextLine()
	}
}

func (p *Parser) processInlineFormatting(text string) string {

	// 1. Inline code (`code`) - must be first to prevent other formatting inside code
	text = inlineCodeRegex.ReplaceAllString(text, "<code>$1</code>")

	// 2. Bold + Italic (***text***)
	text = boldItalicRegex.ReplaceAllString(text, "<strong><em>$1</em></strong>")

	// 3. Bold (**text**)
	text = boldRegex.ReplaceAllString(text, "<strong>$1</strong>")

	// 4. Italic with stars (*text*)
	text = italicStarRegex.ReplaceAllString(text, "<em>$1</em>")

	// 5. Italic with underscores (_text_)
	text = italicUnderRegex.ReplaceAllString(text, "<em>$1</em>")

	return text
}

func (p *Parser) parseBlock() (Node, error) {
	line := strings.TrimSpace(p.current)

	// Check for heading
	if match := headingRegex.FindStringSubmatch(line); match != nil {
		level := len(match[1])
		text := match[2]
		text = p.processInlineFormatting(text)
		return &Heading{Level: level, Text: text}, nil
	}

	if match := blockquoteRegex.FindStringSubmatch(line); match != nil {
		return p.parseBlockquote()
	}

	// Check for code block
	if match := codeBlockRegex.FindStringSubmatch(line); match != nil {
		return p.parseCodeBlock(match[1])
	}

	// Check for list
	if listItemRegex.MatchString(line) {
		return p.parseList()
	}

	// Default to paragraph
	return p.parseParagraph()
}

func (p *Parser) parseBlockquote() (*Blockquote, error) {
	var lines []string

	for {
		match := blockquoteRegex.FindStringSubmatch(p.current)
		if match == nil {
			break
		}
		// Process inline formatting in blockquote content
		content := p.processInlineFormatting(match[1])
		lines = append(lines, content)

		if !p.nextLine() {
			break
		}

		// Check if next line is still part of blockquote or empty line within blockquote
		if !blockquoteRegex.MatchString(p.current) && strings.TrimSpace(p.current) != "" {
			break
		}

		// Handle empty lines within blockquote
		if strings.TrimSpace(p.current) == "" {
			lines = append(lines, "")
			if !p.nextLine() {
				break
			}
		}
	}

	return &Blockquote{Lines: lines}, nil
}

func (p *Parser) parseCodeBlock(language string) (*CodeBlock, error) {
	var lines []string

	for p.nextLine() {
		if strings.HasPrefix(p.current, "```") {
			break
		}
		lines = append(lines, p.current)
	}

	return &CodeBlock{
		Lang:   language,
		Source: lines,
	}, nil
}

// parseList parses an unordered list
func (p *Parser) parseList() (*List, error) {
	var items []string

	for {
		match := listItemRegex.FindStringSubmatch(p.current)
		if match == nil {
			break
		}
		item := p.processInlineFormatting(match[1])
		items = append(items, item)

		if !p.nextLine() {
			break
		}

		// Check if next line is still a list item
		if !listItemRegex.MatchString(p.current) {
			break
		}
	}

	return &List{Items: items}, nil
}

// parseParagraph parses a paragraph (multiple lines until empty line or different block type)
func (p *Parser) parseParagraph() (*Paragraph, error) {
	var lines []string

	for {
		if p.current == "" {
			break
		}

		// Check if this line starts a different block type
		if headingRegex.MatchString(p.current) ||
			codeBlockRegex.MatchString(p.current) ||
			listItemRegex.MatchString(p.current) {
			break
		}

		processedLine := p.processInlineFormatting(strings.TrimSpace(p.current))
		lines = append(lines, processedLine)

		if !p.nextLine() {
			break
		}
	}

	if len(lines) == 0 {
		return nil, nil
	}

	return &Paragraph{Lines: lines}, nil
}

// ParseMarkdown is a convenience function to parse markdown from a reader
func ParseMarkdown(r io.Reader) (*Document, error) {
	parser := NewParser(r)
	return parser.ParserDocument()
}

func (p *Parser) ParserDocument() (*Document, error) {
	doc := &Document{
		Blocks: make([]Node, 0),
	}

	for p.nextLine() {
		p.skipEmptyLines()
		if p.current == "" {
			break
		}
		node, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		if node != nil {
			doc.Blocks = append(doc.Blocks, node)
		}
	}

	return doc, nil
}
