package service

import (
	"context"
	"desafio-clima-cep/internal/client"
)

type WeatherResult struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherService interface {
	GetWeatherByCEP(ctx context.Context, cep string) (*WeatherResult, error)
}

type DefaultWeatherService struct {
	ViaCEPClient     client.ViaCEPClient
	WeatherAPIClient client.WeatherAPIClient
}

func NewDefaultWeatherService(viaCEP client.ViaCEPClient, weatherAPI client.WeatherAPIClient) *DefaultWeatherService {
	return &DefaultWeatherService{
		ViaCEPClient:     viaCEP,
		WeatherAPIClient: weatherAPI,
	}
}

func (s *DefaultWeatherService) GetWeatherByCEP(ctx context.Context, cep string) (*WeatherResult, error) {
	location, err := s.ViaCEPClient.GetLocation(ctx, cep)
	if err != nil {
		return nil, err
	}

	tempC, err := s.WeatherAPIClient.GetTemperature(ctx, location.Localidade)
	if err != nil {
		return nil, err
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	return &WeatherResult{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
