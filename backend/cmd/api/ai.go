// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"
// )

// func summarizeWithGroq(content string) (string, error) {
// 	apiKey := os.Getenv("GROQ_API_KEY")
// 	if apiKey == "" {
// 		return "", fmt.Errorf("missing GROQ_API_KEY environment variable")
// 	}

// 	prompt := `You are a news summarization assistant.

// Summarize the following news articles in no more than 50 words.
// Focus only on the key developments.
// Be factual, concise, and clear.
// Do not add extra information or opinions.

// Articles:
// ` + content

// 	payload := map[string]any{
// 		"model": "llama-3.3-70b-versatile",
// 		"messages": []map[string]string{
// 			{
// 				"role":    "system",
// 				"content": "You summarize news clearly and accurately.",
// 			},
// 			{
// 				"role":    "user",
// 				"content": prompt,
// 			},
// 		},
// 		"temperature": 0.2,
// 	}

// 	bodyBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal Groq request")
// 	}

// 	req, err := http.NewRequest(
// 		http.MethodPost,
// 		"https://api.groq.com/openai/v1/chat/completions",
// 		bytes.NewBuffer(bodyBytes),
// 	)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create Groq request")
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	client := &http.Client{Timeout: 20 * time.Second}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to call Groq API")
// 	}
// 	defer res.Body.Close()

// 	respBody, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read Groq response")
// 	}

// 	if res.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("Groq returned status %d: %s", res.StatusCode, string(respBody))
// 	}

// 	var parsed struct {
// 		Choices []struct {
// 			Message struct {
// 				Content string `json:"content"`
// 			} `json:"message"`
// 		} `json:"choices"`
// 	}

// 	if err := json.Unmarshal(respBody, &parsed); err != nil {
// 		return "", fmt.Errorf("failed to parse Groq response: %s", string(respBody))
// 	}

// 	if len(parsed.Choices) == 0 {
// 		return "", fmt.Errorf("Groq response contained no choices")
// 	}

//		return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
//	}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func summarizeWithGroq(content string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing GROQ_API_KEY environment variable")
	}

	prompt := `Summarize the following news article in one clear paragraph.

Rules:
- Write between 60 and 70 words
- No bullet points
- No numbering
- No intro text
- Be factual and clear
- Cover the main development and key context
- Only use information from the article

Article:
` + content

	payload := map[string]any{
		"model": "llama-3.3-70b-versatile",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You summarize one news article into one plain-text paragraph between 60 and 70 words.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.2,
		"max_tokens":  150,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Groq request")
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create Groq request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 20 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Groq API")
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Groq response")
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Groq returned status %d: %s", res.StatusCode, string(respBody))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse Groq response: %s", string(respBody))
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("Groq response contained no choices")
	}

	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}
