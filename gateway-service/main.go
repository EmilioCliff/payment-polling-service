package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Printf("Failed to start new server instance to db: %s", err)
		return
	}

	log.Printf("Starting Gateway Server at port: %s", os.Getenv("PORT"))
	server.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
}
