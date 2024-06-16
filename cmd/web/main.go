package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type app struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":5123", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &app{
		logger: logger,
	}

	logger.Info(fmt.Sprintf("Listening on http://localhost%s", *addr))
	err := http.ListenAndServe(*addr, app.routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
