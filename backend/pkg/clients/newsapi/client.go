package newsapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	APIKey string
	BaseURL string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		BaseURL: "https://newsapi.org/v2",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Response structs
type Article struct {
	Source      string `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Content     string `json:"content"`
}

type NewsAPIResponse struct {
	Status  string    `json:"status"`
	TotalResults int  `json:"totalResults"`
	Articles []Article `json:"articles"`
}

// Fetch top headlines for a topic
func (c *Client) GetTopHeadlines(topic string) ([]Article, error) {
	url := fmt.Sprintf("%s/top-headlines?q=%s&apiKey=%s", c.BaseURL, topic, c.APIKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("NewsAPI error: %s", string(body))
	}

	var result NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return result.Articles, nil
}