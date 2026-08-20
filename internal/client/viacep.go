package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Erro        bool   `json:"erro,omitempty"`
}

type ViaCEPClient interface {
	GetLocation(ctx context.Context, cep string) (*ViaCEPResponse, error)
}

type HTTPViaCEPClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPViaCEPClient(baseURL string) *HTTPViaCEPClient {
	if baseURL == "" {
		baseURL = "https://viacep.com.br"
	}
	return &HTTPViaCEPClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPViaCEPClient) GetLocation(ctx context.Context, cep string) (*ViaCEPResponse, error) {
	url := fmt.Sprintf("%s/ws/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("zipcode not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("viacep api returned status: %d", resp.StatusCode)
	}

	var result ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Erro || result.Localidade == "" {
		return nil, fmt.Errorf("zipcode not found")
	}

	return &result, nil
}
