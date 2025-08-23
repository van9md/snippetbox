package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/van9md/snippetbox/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	logger   *slog.Logger
	cfg      config
	snippets *models.SnippetModel
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
		cfg:    config{},
	}

	flag.StringVar(&app.cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&app.cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL datasource name")
	flag.Parse()

	logger.Info("starting server", slog.String("addr", app.cfg.addr))

	db, err := openDB(app.cfg.dsn)
	defer db.Close()
	app.snippets = &models.SnippetModel{DB: db}

	err = http.ListenAndServe(app.cfg.addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
