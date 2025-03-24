package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type ExchangeRateResponse struct {
    Dolar string `json:"Dólar"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Error making request:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        fmt.Println("Error response from server:", string(body))
        return
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return
    }

    var exchangeRate ExchangeRateResponse
    if err := json.Unmarshal(body, &exchangeRate); err != nil {
        fmt.Println("Error unmarshaling response:", err)
        return
    }

    fileContent := fmt.Sprintf("Dólar: %s", exchangeRate.Dolar)
    if err := os.WriteFile("cotacao.txt", []byte(fileContent), os.ModePerm); err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    fmt.Println("Saved exchange rate to cotacao.txt")
}