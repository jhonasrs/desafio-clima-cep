package handler

import (
	"context"
	"desafio-clima-cep/internal/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockWeatherService struct {
	result *service.WeatherResult
	err    error
}

func (m *mockWeatherService) GetWeatherByCEP(ctx context.Context, cep string) (*service.WeatherResult, error) {
	return m.result, m.err
}

func TestWeatherHandler_Success(t *testing.T) {
	svc := &mockWeatherService{
		result: &service.WeatherResult{
			TempC: 30.0,
			TempF: 86.0,
			TempK: 303.0,
		},
	}
	h := NewWeatherHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/cep/01153000", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	expectedBody := `{"temp_C":30,"temp_F":86,"temp_K":303}`
	// Note: json encoder format might have newline or compact. Let's check body.
	if rec.Body.String() != expectedBody+"\n" && rec.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	svc := &mockWeatherService{}
	h := NewWeatherHandler(svc)

	// Less than 8 digits or non-numeric
	req := httptest.NewRequest(http.MethodGet, "/cep/12345", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rec.Code)
	}

	expectedBody := `{"message":"invalid zipcode"}`
	if rec.Body.String() != expectedBody+"\n" && rec.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}

func TestWeatherHandler_NotFoundCEP(t *testing.T) {
	svc := &mockWeatherService{
		err: errors.New("zipcode not found"),
	}
	h := NewWeatherHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/cep/99999999", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	expectedBody := `{"message":"can not find zipcode"}`
	if rec.Body.String() != expectedBody+"\n" && rec.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}
