package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/api"
)

func main() {
	server := api.NewServer()

	log.Printf("Starting Gateway Server at port: %s", os.Getenv("PORT"))
	server.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
}
