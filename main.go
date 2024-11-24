package main

import (
	"ecommerce/wire"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, err := wire.ProvideLogger()
	if err != nil {
		log.Panicf("Failed to initialize logger: %v", err)
	}

	r, err := wire.InitializeRouterHandler()
	if err != nil {
		logger.Panic("Failed to initialize router", zap.Error(err))
	}
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))

}
