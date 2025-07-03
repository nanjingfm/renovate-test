// Package netutil provides network utility functions using golang.org/x/net
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HttpClient wraps http.Client with context support
type HttpClient struct {
	client *http.Client
}

// NewHttpClient creates a new HttpClient with timeout
func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get performs an HTTP GET request with context
func (c *HttpClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Use context from x/net/context
	req = req.WithContext(ctx)

	return c.client.Do(req)
}

// HtmlParser provides HTML parsing functionality
type HtmlParser struct{}

// NewHtmlParser creates a new HtmlParser instance
func NewHtmlParser() *HtmlParser {
	return &HtmlParser{}
}

// ExtractTitle extracts the title from HTML content
func (p *HtmlParser) ExtractTitle(htmlContent string) (string, error) {
	reader := strings.NewReader(htmlContent)
	doc, err := html.Parse(reader)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	title := p.findTitle(doc)
	return title, nil
}

// findTitle recursively searches for the title element
func (p *HtmlParser) findTitle(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "title" {
		if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
			return node.FirstChild.Data
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if title := p.findTitle(child); title != "" {
			return title
		}
	}

	return ""
}

// GetPageTitle fetches a webpage and extracts its title
func GetPageTitle(ctx context.Context, url string) (string, error) {
	client := NewHttpClient(30 * time.Second)
	parser := NewHtmlParser()

	// Fetch the webpage
	resp, err := client.Get(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract title from HTML
	title, err := parser.ExtractTitle(string(body))
	if err != nil {
		return "", fmt.Errorf("failed to extract title: %w", err)
	}

	return title, nil
}
