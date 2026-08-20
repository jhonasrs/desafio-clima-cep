package service

import (
	"context"
	"desafio-clima-cep/internal/client"
	"errors"
	"testing"
)

type mockViaCEPClient struct {
	resp *client.ViaCEPResponse
	err  error
}

func (m *mockViaCEPClient) GetLocation(ctx context.Context, cep string) (*client.ViaCEPResponse, error) {
	return m.resp, m.err
}

type mockWeatherAPIClient struct {
	temp float64
	err  error
}

func (m *mockWeatherAPIClient) GetTemperature(ctx context.Context, city string) (float64, error) {
	return m.temp, m.err
}

func TestGetWeatherByCEP_Success(t *testing.T) {
	viaCEP := &mockViaCEPClient{
		resp: &client.ViaCEPResponse{
			CEP:        "01153000",
			Localidade: "São Paulo",
		},
	}
	weatherAPI := &mockWeatherAPIClient{
		temp: 25.0,
	}

	svc := NewDefaultWeatherService(viaCEP, weatherAPI)
	res, err := svc.GetWeatherByCEP(context.Background(), "01153000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.TempC != 25.0 {
		t.Errorf("expected TempC 25.0, got %f", res.TempC)
	}
	// TempF = 25 * 1.8 + 32 = 45 + 32 = 77
	if res.TempF != 77.0 {
		t.Errorf("expected TempF 77.0, got %f", res.TempF)
	}
	// TempK = 25 + 273 = 298
	if res.TempK != 298.0 {
		t.Errorf("expected TempK 298.0, got %f", res.TempK)
	}
}

func TestGetWeatherByCEP_ViaCEPError(t *testing.T) {
	viaCEP := &mockViaCEPClient{
		err: errors.New("zipcode not found"),
	}
	weatherAPI := &mockWeatherAPIClient{
		temp: 25.0,
	}

	svc := NewDefaultWeatherService(viaCEP, weatherAPI)
	_, err := svc.GetWeatherByCEP(context.Background(), "99999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
