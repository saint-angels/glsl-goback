package main

import (
	"flag"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"saint-angels/shaderbox/pkg/models"
	"saint-angels/shaderbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"

	"saint-angels/shaderbox/pkg/renderer"
)

type contextKey string
const contextKeyIsAuthenticated = contextKey("isAuthenticated")


type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	artworks interface {
		Insert() (int, error)
		GetArtForRender() (*models.Artwork, error)
	}
	collector *renderer.Collector
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/shaderbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	//Get pool of database connections
    db, err := sql.Open("mysql", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
    if err = db.Ping(); err != nil {
		errorLog.Fatal(err)
    }
	defer db.Close()

	artworkmodel := &mysql.ArtworkModel{DB: db}
	app := &application {
		infoLog: infoLog,
		errorLog: errorLog,
		artworks: artworkmodel,
	}

	//Start the render worker pool
	collector := renderer.StartDispatcher(1, artworkmodel)
	defer func() {
		collector.End <- true
	}()

	srv := &http.Server {
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,

	}

	infoLog.Printf("starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
