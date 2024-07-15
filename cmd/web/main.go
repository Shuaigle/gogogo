package main

import (
	"flag"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"pinkamkak.com/web/internal/models"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pinkamkak@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// db
	db, dbErr := openDB(*dsn)
	if dbErr != nil {
		logger.Error(dbErr.Error())
		os.Exit(1)
	}
	defer db.Close()

	// template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	mux := app.routes()
	logger.Info("Starting server", zap.String("address", *addr))
	httpErr := http.ListenAndServe(*addr, mux)
	if httpErr != nil {
		logger.Fatal("Server failed to start", zap.Error(httpErr))
		os.Exit(1)
	}
}
