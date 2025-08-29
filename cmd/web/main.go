package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/van9md/snippetbox/internal/models"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	logger         *slog.Logger
	cfg            config
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		logger:        logger,
		cfg:           config{},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	flag.StringVar(&app.cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&app.cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL datasource name")
	flag.Parse()

	logger.Info("starting server", slog.String("addr", app.cfg.addr))

	db, err := openDB(app.cfg.dsn)
	defer db.Close()
	app.snippets = &models.SnippetModel{DB: db}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	app.sessionManager = sessionManager

	srv := &http.Server{
		Addr:     app.cfg.addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err = srv.ListenAndServe()
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
