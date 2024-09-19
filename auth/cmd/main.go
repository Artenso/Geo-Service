package main

import (
	"context"
	"log"

	"github.com/Artenso/auth-service/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %s", err)
	}

	if err := a.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %s", err)
	}
}
