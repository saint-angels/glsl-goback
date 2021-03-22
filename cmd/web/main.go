package main

import (
	"flag"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"saint-angels/glsl-goback/pkg/models"
	"saint-angels/glsl-goback/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type contextKey string
const contextKeyIsAuthenticated = contextKey("isAuthenticated")


type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	artworks interface {
		Insert() (int, error)
		GetOldestUnrendered() (*models.Artwork, error)
	}
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

	app := &application {
		infoLog: infoLog,
		errorLog: errorLog,
		artworks: &mysql.ArtworkModel{DB: db},
	}

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
