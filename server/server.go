package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luuan11/client-server/server/database"
	"github.com/luuan11/client-server/server/handler"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
        fmt.Println("Failed to initialize database")
	}

	defer database.DB.Close()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Use /cotacao to get the current exchange rate."))
    })

	http.HandleFunc("/cotacao", handler.ExchangeHandler)

	log.Default().Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}