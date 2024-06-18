package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/amirintech/snippet-box/internal/models"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type app struct {
	logger        *slog.Logger
	snippetModel  *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":5123", "HTTP network address")
	dsn := flag.String("dsn", "", "database connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := connectToDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &app{
		logger:        logger,
		snippetModel:  &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info(fmt.Sprintf("Listening on http://localhost%s", *addr))
	err = http.ListenAndServe(*addr, app.routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func connectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
