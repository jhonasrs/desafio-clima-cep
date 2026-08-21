# Serviço de Clima por CEP (CEP Weather Service)

Uma aplicação desenvolvida em Go que recebe um CEP brasileiro, valida o formato, consulta o serviço [`ViaCEP`](internal/client/viacep.go:1) para identificar a cidade correspondente, obtém a temperatura atual na [`WeatherAPI`](internal/client/weatherapi.go:1), realiza a conversão para Fahrenheit e Kelvin, e retorna uma resposta JSON formatada. O sistema também está preparado para containerização com [`Docker`](Dockerfile:1) ou [`Podman`](Dockerfile:1).

---

## 📋 Requisitos

- **Go**: Versão 1.25.4 ou superior ([`go.mod`](go.mod:3))
- **Docker** ou **Podman** (opcional para execução containerizada)

---

## 🧪 Executando os Testes

Para executar todos os testes unitários e de integração (`go test ./...`), execute no terminal:

```bash
go test -v ./...
```

---

## 🚀 Executando Localmente

### 1. Executando com Go
Execute o servidor [`main.go`](cmd/server/main.go:1):

```bash
go run cmd/server/main.go
```

O servidor iniciará na porta `8080` (configurável via variável de ambiente `PORT`).

### 2. Executando com Docker ou Podman
Construa e execute a imagem [`Dockerfile`](Dockerfile:1) usando Docker:

```bash
docker build -t desafio-clima-cep .
docker run -p 8080:8080 desafio-clima-cep
```

Ou usando Podman:

```bash
podman build -t desafio-clima-cep .
podman run -p 8080:8080 desafio-clima-cep
```

---

## 🔌 Exemplos de Uso da API

### Sucesso (`200 OK`)
**Requisição:**
```bash
curl -i http://localhost:8080/cep/01153000
curl -i https://desafio-clima-cep-1002536084691.us-east1.run.app/cep/01153000
```

**Resposta:**
```json
{
  "temp_C": 25.0,
  "temp_F": 77.0,
  "temp_K": 298.0
}
```

### Formato de CEP Inválido (`422 Unprocessable Entity`)
**Requisição:**
```bash
curl -i http://localhost:8080/cep/12345
curl -i https://desafio-clima-cep-1002536084691.us-east1.run.app/cep/12345
```

**Resposta:**
```json
{
  "message": "invalid zipcode"
}
```

### CEP Não Encontrado (`404 Not Found`)
**Requisição:**
```bash
curl -i http://localhost:8080/cep/99999999
curl -i https://desafio-clima-cep-1002536084691.us-east1.run.app/cep/99999999
```

**Resposta:**
```json
{
  "message": "can not find zipcode"
}
```
