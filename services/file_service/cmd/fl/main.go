package main

import (
	"context"
	"flag"

	"github.com/alexey-dobry/fileshare/services/file_service/internal/app"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/config"
)

func main() {
	flag.Parse()

	cfg := config.MustLoad()

	application := app.New(cfg)

	application.Run(context.Background())
}
