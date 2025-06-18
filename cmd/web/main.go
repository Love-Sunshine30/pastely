package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"al.imran.pastely/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

// defining application struct to hold applicatiom-wide dependencies
type application struct {
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	snippets       *models.SnippetModel
	user           *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// defining a flag named addr with default value ":4000"
	addr := flag.String("addr", ":4000", "HTTP network string")

	// definig a new command-line flag for the mySQL DSN string
	dsn := flag.String("dsn", "web:whoami7@/pastely?parseTime=true", "mySQL data source name")

	flag.Parse()

	// creating two new Logger. one for INFO and another for ERROR message

	// info logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// error logger
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// creating a connection pool
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// the connection pool will be closed before the main function exits
	defer db.Close()

	// creating a new instances of templateCache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	// initializing formDecoder instance
	formDecoder := form.NewDecoder()

	// creating a new session maganger and connecting it to our mysql sever
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// cookie will only be sent over https connection
	sessionManager.Cookie.Secure = true

	//creating an instance of application struct
	app := &application{
		errorLogger:    errorLog,
		infoLogger:     infoLog,
		snippets:       &models.SnippetModel{DB: db},
		user:           &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// creating a http server struct to contain everything the server needs to run
	// it would contain network address, handler, custom logger

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(), // this contains all the routes
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Strating server on port %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// openDB() function returns a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
