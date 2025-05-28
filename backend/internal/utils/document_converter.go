package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DocumentConverter provides document format conversion
type DocumentConverter struct{}

// NewDocumentConverter creates a new document converter
func NewDocumentConverter() *DocumentConverter {
	return &DocumentConverter{}
}

// ConvertToMarkdown converts document to markdown format
func (dc *DocumentConverter) ConvertToMarkdown(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Simple conversion based on input type
	ext := strings.ToLower(filepath.Ext(inputPath))
	var markdown string

	switch ext {
	case ".txt":
		markdown = dc.convertTextToMarkdown(string(content))
	case ".html", ".htm":
		markdown = dc.convertHTMLToMarkdown(string(content))
	default:
		markdown = fmt.Sprintf("# %s\n\n%s", filepath.Base(inputPath), string(content))
	}

	return os.WriteFile(outputPath, []byte(markdown), 0644)
}

// ConvertToHTML converts document to HTML format
func (dc *DocumentConverter) ConvertToHTML(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	var html string

	switch ext {
	case ".md", ".markdown":
		html = dc.convertMarkdownToHTML(string(content))
	case ".txt":
		html = dc.convertTextToHTML(string(content))
	default:
		html = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <meta charset="UTF-8">
</head>
<body>
    <pre>%s</pre>
</body>
</html>`, filepath.Base(inputPath), string(content))
	}

	return os.WriteFile(outputPath, []byte(html), 0644)
}

// ConvertToPlainText converts document to plain text
func (dc *DocumentConverter) ConvertToPlainText(inputPath, outputPath string) error {
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	var plainText string

	switch ext {
	case ".html", ".htm":
		plainText = StripHTML(string(content))
	case ".md", ".markdown":
		plainText = dc.convertMarkdownToText(string(content))
	default:
		plainText = string(content)
	}

	return os.WriteFile(outputPath, []byte(plainText), 0644)
}

// Helper methods for conversion

func (dc *DocumentConverter) convertTextToMarkdown(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	result.WriteString("# Document\n\n")

	inCodeBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect code blocks (lines starting with spaces)
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			if !inCodeBlock {
				result.WriteString("```\n")
				inCodeBlock = true
			}
			result.WriteString(line + "\n")
		} else {
			if inCodeBlock {
				result.WriteString("```\n\n")
				inCodeBlock = false
			}

			// Detect potential headers
			if len(trimmed) > 0 && len(trimmed) < 100 && dc.isPotentialHeader(trimmed) {
				result.WriteString("## " + trimmed + "\n\n")
			} else {
				result.WriteString(line + "\n")
			}
		}
	}

	if inCodeBlock {
		result.WriteString("```\n")
	}

	return result.String()
}

func (dc *DocumentConverter) convertMarkdownToHTML(markdown string) string {
	// Basic markdown to HTML conversion
	html := strings.ReplaceAll(markdown, "\n", "<br>\n")

	// Headers
	html = regexp.MustCompile(`(?m)^### (.+)$`).ReplaceAllString(html, "<h3>$1</h3>")
	html = regexp.MustCompile(`(?m)^## (.+)$`).ReplaceAllString(html, "<h2>$1</h2>")
	html = regexp.MustCompile(`(?m)^# (.+)$`).ReplaceAllString(html, "<h1>$1</h1>")

	// Bold and italic
	html = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(html, "<strong>$1</strong>")
	html = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(html, "<em>$1</em>")

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Converted Document</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1, h2, h3 { color: #333; }
        pre { background: #f4f4f4; padding: 10px; }
    </style>
</head>
<body>
%s
</body>
</html>`, html)
}

func (dc *DocumentConverter) convertTextToHTML(text string) string {
	// Escape HTML chars and convert newlines
	html := strings.ReplaceAll(text, "&", "&amp;")
	html = strings.ReplaceAll(html, "<", "&lt;")
	html = strings.ReplaceAll(html, ">", "&gt;")
	html = strings.ReplaceAll(html, "\n", "<br>\n")

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Text Document</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: monospace; margin: 40px; white-space: pre-wrap; }
    </style>
</head>
<body>
%s
</body>
</html>`, html)
}

func (dc *DocumentConverter) convertMarkdownToText(markdown string) string {
	// Remove markdown formatting
	text := markdown

	// Remove headers
	text = regexp.MustCompile(`(?m)^#+\s*`).ReplaceAllString(text, "")

	// Remove formatting
	text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile("`(.+?)`").ReplaceAllString(text, "$1")

	return text
}

func (dc *DocumentConverter) convertHTMLToMarkdown(htmlContent string) string {
	// Basic HTML to Markdown conversion
	content := htmlContent

	// Convert headers
	content = regexp.MustCompile(`<h1[^>]*>(.*?)</h1>`).ReplaceAllString(content, "# $1")
	content = regexp.MustCompile(`<h2[^>]*>(.*?)</h2>`).ReplaceAllString(content, "## $1")
	content = regexp.MustCompile(`<h3[^>]*>(.*?)</h3>`).ReplaceAllString(content, "### $1")
	content = regexp.MustCompile(`<h4[^>]*>(.*?)</h4>`).ReplaceAllString(content, "#### $1")
	content = regexp.MustCompile(`<h5[^>]*>(.*?)</h5>`).ReplaceAllString(content, "##### $1")
	content = regexp.MustCompile(`<h6[^>]*>(.*?)</h6>`).ReplaceAllString(content, "###### $1")

	// Convert formatting
	content = regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(content, "**$1**")
	content = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(content, "**$1**")
	content = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(content, "*$1*")
	content = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(content, "*$1*")
	content = regexp.MustCompile(`<code[^>]*>(.*?)</code>`).ReplaceAllString(content, "`$1`")

	// Convert links
	content = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>(.*?)</a>`).ReplaceAllString(content, "[$2]($1)")

	// Convert images
	content = regexp.MustCompile(`<img[^>]*src="([^"]*)"[^>]*alt="([^"]*)"[^>]*/?>`).ReplaceAllString(content, "![$2]($1)")
	content = regexp.MustCompile(`<img[^>]*src="([^"]*)"[^>]*/?>`).ReplaceAllString(content, "![]($1)")

	// Convert paragraphs
	content = regexp.MustCompile(`<p[^>]*>(.*?)</p>`).ReplaceAllString(content, "$1\n\n")

	// Convert line breaks
	content = regexp.MustCompile(`<br\s*/?>|<br>`).ReplaceAllString(content, "\n")

	// Convert lists
	content = regexp.MustCompile(`<ul[^>]*>`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`</ul>`).ReplaceAllString(content, "\n")
	content = regexp.MustCompile(`<ol[^>]*>`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`</ol>`).ReplaceAllString(content, "\n")
	content = regexp.MustCompile(`<li[^>]*>(.*?)</li>`).ReplaceAllString(content, "- $1")

	// Convert code blocks
	content = regexp.MustCompile(`<pre[^>]*><code[^>]*>(.*?)</code></pre>`).ReplaceAllString(content, "```\n$1\n```")
	content = regexp.MustCompile(`<pre[^>]*>(.*?)</pre>`).ReplaceAllString(content, "```\n$1\n```")

	// Remove remaining HTML tags
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")

	// Clean up whitespace
	content = regexp.MustCompile(`\n\s*\n\s*\n`).ReplaceAllString(content, "\n\n")
	content = strings.TrimSpace(content)

	return content
}

func (dc *DocumentConverter) isPotentialHeader(line string) bool {
	// Simple heuristics for header detection
	if len(line) > 100 {
		return false
	}

	// All caps might be a header
	if strings.ToUpper(line) == line && len(line) > 3 {
		return true
	}

	// Numbered sections
	if matched, _ := regexp.MatchString(`^\d+\.?\s+[A-Z]`, line); matched {
		return true
	}

	return false
}
