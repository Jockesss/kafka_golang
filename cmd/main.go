package main

import (
	"context"
	"log"

	"kafka-golang/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("start app error: %v", err)
	}
}
