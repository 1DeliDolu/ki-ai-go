// backend/internal/services/wiki_service.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/1DeliDolu/go_mustAI/local-ai-project/backend/pkg/types"
)

type WikiService struct {
	baseURL string
}

func NewWikiService() *WikiService {
	return &WikiService{
		baseURL: "https://de.wikipedia.org/api/rest_v1",
	}
}

func (s *WikiService) Search(query string) ([]types.WikiResult, error) {
	// Wikipedia search API
	searchURL := fmt.Sprintf("%s/page/summary/%s", s.baseURL, url.QueryEscape(query))

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try search API instead
		return s.searchMultiple(query)
	}

	var result struct {
		Title       string `json:"title"`
		Extract     string `json:"extract"`
		Description string `json:"description"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
		Thumbnail struct {
			Source string `json:"source"`
		} `json:"thumbnail"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return []types.WikiResult{
		{
			Title:       result.Title,
			Extract:     result.Extract,
			Description: result.Description,
			URL:         result.ContentURLs.Desktop.Page,
			Thumbnail:   result.Thumbnail.Source,
		},
	}, nil
}

func (s *WikiService) searchMultiple(query string) ([]types.WikiResult, error) {
	// Use OpenSearch API for multiple results
	searchURL := fmt.Sprintf("https://de.wikipedia.org/w/api.php?action=opensearch&search=%s&limit=5&format=json",
		url.QueryEscape(query))

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResults []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return nil, err
	}

	if len(searchResults) < 4 {
		return []types.WikiResult{}, nil
	}

	titles, ok := searchResults[1].([]interface{})
	if !ok {
		return []types.WikiResult{}, nil
	}

	descriptions, ok := searchResults[2].([]interface{})
	if !ok {
		descriptions = make([]interface{}, len(titles))
	}

	urls, ok := searchResults[3].([]interface{})
	if !ok {
		return []types.WikiResult{}, nil
	}

	var results []types.WikiResult
	for i, title := range titles {
		if i < len(descriptions) && i < len(urls) {
			desc := ""
			if descriptions[i] != nil {
				desc = descriptions[i].(string)
			}

			results = append(results, types.WikiResult{
				Title:       title.(string),
				Description: desc,
				Extract:     desc,
				URL:         urls[i].(string),
			})
		}
	}

	return results, nil
}
