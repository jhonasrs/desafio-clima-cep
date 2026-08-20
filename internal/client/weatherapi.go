package client

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//go:embed api-key.txt
var embeddedAPIKey string

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherAPIClient interface {
	GetTemperature(ctx context.Context, city string) (float64, error)
}

type HTTPWeatherAPIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewHTTPWeatherAPIClient(baseURL string) *HTTPWeatherAPIClient {
	if baseURL == "" {
		baseURL = "https://api.weatherapi.com"
	}
	var apiKey string
	if data, err := os.ReadFile("api-key.txt"); err == nil {
		apiKey = strings.TrimSpace(string(data))
	} else if data, err := os.ReadFile("../api-key.txt"); err == nil {
		apiKey = strings.TrimSpace(string(data))
	} else if data, err := os.ReadFile("../../api-key.txt"); err == nil {
		apiKey = strings.TrimSpace(string(data))
	}
	if apiKey == "" {
		apiKey = strings.TrimSpace(embeddedAPIKey)
	}

	return &HTTPWeatherAPIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPWeatherAPIClient) GetTemperature(ctx context.Context, city string) (float64, error) {
	escapedCity := url.QueryEscape(city)
	reqURL := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", c.BaseURL, c.APIKey, escapedCity)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi returned status: %d", resp.StatusCode)
	}

	var result WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Current.TempC, nil
}
