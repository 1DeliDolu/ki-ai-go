package utils

import (
	"html"
	"regexp"
	"strings"
)

// FileInfo represents comprehensive file information

// DetectLanguage provides basic language detection
func DetectLanguage(text string) string {
	// Simple heuristic-based language detection
	text = strings.ToLower(text)

	// English indicators
	englishWords := []string{"the", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"}
	englishCount := 0
	for _, word := range englishWords {
		englishCount += strings.Count(text, " "+word+" ")
	}

	// German indicators
	germanWords := []string{"der", "die", "das", "und", "oder", "aber", "in", "auf", "mit", "von", "zu", "für"}
	germanCount := 0
	for _, word := range germanWords {
		germanCount += strings.Count(text, " "+word+" ")
	}

	// Turkish indicators
	turkishWords := []string{"ve", "veya", "ama", "ile", "den", "dan", "için", "gibi", "kadar", "daha"}
	turkishCount := 0
	for _, word := range turkishWords {
		turkishCount += strings.Count(text, " "+word+" ")
	}

	if englishCount > germanCount && englishCount > turkishCount {
		return "en"
	} else if germanCount > turkishCount {
		return "de"
	} else if turkishCount > 0 {
		return "tr"
	}

	return "unknown"
}

// CalculateComplexityScore calculates text complexity (0-100)
func CalculateComplexityScore(text string) int {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0
	}

	sentences := strings.FieldsFunc(text, func(c rune) bool {
		return c == '.' || c == '!' || c == '?'
	})

	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))

	// Count long words (>6 characters)
	longWords := 0
	totalChars := 0
	for _, word := range words {
		if len(word) > 6 {
			longWords++
		}
		totalChars += len(word)
	}

	avgWordLength := float64(totalChars) / float64(len(words))
	longWordRatio := float64(longWords) / float64(len(words))

	// Simple complexity formula
	complexity := int((avgWordsPerSentence * 2) + (avgWordLength * 10) + (longWordRatio * 50))

	if complexity > 100 {
		complexity = 100
	}

	return complexity
}

// StripHTML removes HTML tags from text
func StripHTML(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	stripped := re.ReplaceAllString(content, "")
	return html.UnescapeString(stripped)
}

// CountWords counts words in text
func CountWords(text string) int {
	words := strings.Fields(strings.TrimSpace(text))
	return len(words)
}

// ExtractLinks extracts URLs from text
func ExtractLinks(text string) []string {
	urlPattern := `https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`
	re := regexp.MustCompile(urlPattern)
	return re.FindAllString(text, -1)
}

// TruncateString truncates string to specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// CleanText performs basic text cleaning
func CleanText(text string) string {
	// Remove excessive whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	return text
}
