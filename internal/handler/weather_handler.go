package handler

import (
	"desafio-clima-cep/internal/service"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

type WeatherHandler struct {
	WeatherService service.WeatherService
	CEPRegex       *regexp.Regexp
}

func NewWeatherHandler(ws service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		WeatherService: ws,
		CEPRegex:       regexp.MustCompile(`^\d{8}$`),
	}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract CEP from URL path, e.g., /cep/{cep} or /cep/01153000
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Expected path format: ["cep", "{cep}"] or similar, or just handling /cep/{cep}
	var cep string
	if len(parts) >= 2 && parts[0] == "cep" {
		cep = parts[1]
	} else if len(parts) == 1 && parts[0] != "" {
		cep = parts[0]
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message":"invalid zipcode"}`))
		return
	}

	if !h.CEPRegex.MatchString(cep) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message":"invalid zipcode"}`))
		return
	}

	weather, err := h.WeatherService.GetWeatherByCEP(r.Context(), cep)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"can not find zipcode"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}
