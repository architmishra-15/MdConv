package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// GenerateHighlightDiv builds an HTML snippet for a single code block,
// scoped to a unique container so that Highlight.js CSS and JS apply only to it.
// - lang: the language after ``` (e.g. "go" or "python").
// - source: raw code lines.
// - style: name of the highlight.js theme (without .min.css).
// Assumes files exist under ./lang_highlight/styles/{style}.min.css and
// ./lang_highlight/highlight.min.js and ./lang_highlight/languages/{lang}.min.js.

func GenerateHighlightDiv(lang string, source []string, style string, id string) (template.HTML, error) {
	// Read and scope the CSS
	cssPath := filepath.Join("lang_highlight", "styles", style+".min.css")
	cssBytes, err := ioutil.ReadFile(cssPath)
	if err != nil {
		// fallback to github-dark theme
		cssBytes, _ = ioutil.ReadFile(filepath.Join("lang_highlight", "styles", "github-dark.min.css"))
	}
	rawCSS := string(cssBytes)
	// Prefix every selector with #<id> to scope styles
	lines := strings.Split(rawCSS, "\n")
	var scoped []string
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "@") {
			// Keep @import or @keyframes as is
			scoped = append(scoped, line)
		} else if idx := strings.Index(line, "{"); idx != -1 {
			// prefix selectors before '{'
			sel := line[:idx]
			rest := line[idx:]
			prefixed := fmt.Sprintf("#%s %s%s", id, sel, rest)
			scoped = append(scoped, prefixed)
		} else {
			scoped = append(scoped, line)
		}
	}
	// scopedCSS := strings.Join(scoped, "\n")

	// Build the HTML snippet
	code := template.HTMLEscapeString(strings.Join(source, "\n"))

	html := fmt.Sprintf(`
<style>%s</style>
<div id="%s" class="hljs-container">
  <pre><code class="language-%s">%s</code></pre>
</div>
<script>
  document.addEventListener("DOMContentLoaded", () => {
    const el = document.querySelector("#%s code");
    if (window.hljs) hljs.highlightElement(el);
  });
</script>`, strings.Join(scoped, "\n"), id, lang, code, id)

	return template.HTML(html), nil
}
