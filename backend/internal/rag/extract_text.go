package rag

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenthenguyen/docx"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"golang.org/x/net/html"
)

func ExtractText(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pdf":
		return extractPDF(path)
	case ".txt":
		return extractTXT(path)
	case ".doc", ".docx":
		return extractDOCX(path)
	case ".md":
		return extractTXT(path) // Markdown da düz metin gibi işlenebilir
	case ".html", ".htm":
		return extractHTML(path)
	default:
		return "", ErrUnsupportedFile
	}
}

func extractPDF(path string) (string, error) {
	reader, _, err := model.NewPdfReaderFromFile(path, nil)
	if err != nil {
		return "", err
	}
	numPages, _ := reader.GetNumPages()
	var allText string
	for i := 1; i <= numPages; i++ {
		page, _ := reader.GetPage(i)
		ex, _ := extractor.New(page)
		text, _ := ex.ExtractText()
		allText += text + "\n"
	}
	return allText, nil
}

func extractTXT(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}

func extractDOCX(path string) (string, error) {
	r, err := docx.ReadDocxFile(path)
	if err != nil {
		return "", err
	}
	defer r.Close()
	doc := r.Editable()
	return doc.GetContent(), nil
}

func extractHTML(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	doc, err := html.Parse(f)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data + " ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)
	return sb.String(), nil
}

var ErrUnsupportedFile = errors.New("unsupported file type")
