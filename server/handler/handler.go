package handler

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/luuan11/client-server/server/database"
    "github.com/luuan11/client-server/server/models"
)

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
    reqCtx, reqCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
    defer reqCancel()

    resp, err := fetchUSD(reqCtx)
    if err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Timeout exceeded while fetching exchange rate", http.StatusRequestTimeout)
        } else {
            http.Error(w, "Error fetching exchange rate: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
    defer dbCancel()

    if err := database.InsertNewExchangeRate(dbCtx, resp.Bid); err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Timeout exceeded while inserting exchange rate", http.StatusRequestTimeout)
        } else {
            http.Error(w, "Error inserting exchange rate: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    jsonResp := map[string]string{"Dólar": resp.Bid}
    if err := json.NewEncoder(w).Encode(jsonResp); err != nil {
        http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
        return
    }
}

func fetchUSD(ctx context.Context) (*models.ExchangeApiResponse, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
    if err != nil {
        fmt.Println("Failed to create request")
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Failed to send request")
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Failed to read response body")
        return nil, err
    }

    var data models.ExchangeApiResponse
    if err = json.Unmarshal(body, &data); err != nil {
        return nil, err
    }

    return &data, nil
}