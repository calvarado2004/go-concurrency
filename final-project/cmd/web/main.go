package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/calvarado2004/go-concurrency/data"
	"github.com/gomodule/redigo/redis"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8080"

func main() {

	// create web sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	panicLog := log.New(os.Stderr, "PANIC\t", log.Ldate|log.Ltime|log.Lshortfile)

	// set up application config just for logging
	app := Config{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		PanicLog: panicLog,
	}

	// connect to the database
	db := initDB(&app)

	// create channels

	// create wait group
	wg := sync.WaitGroup{}

	// set up application config, including db and wait group
	app = Config{
		Session:       session,
		DB:            db,
		Wait:          &wg,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		PanicLog:      panicLog,
		Models:        data.New(db),
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
	}

	// set up email
	app.Mailer = app.CreateMail()

	go app.listenForMail()

	// listen for signals
	go app.listenForShutdown()

	// listen for errors
	go app.listenForErrors()

	// listen for web connections
	app.serve()

}

func (app *Config) listenForErrors() {
	for {
		select {
		case err := <-app.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.ErrorChanDone:
			return
		}
	}
}

func (app *Config) serve() {
	// serve the web application
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Printf("Starting server on port %s", webPort)

	err := srv.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatal(err)
	}

}

func initDB(app *Config) *sql.DB {

	// connect to the database
	conn := connectToDB(app)
	if conn == nil {
		app.PanicLog.Println("We could not connect to the database after 10 attempts")
	}

	return conn
}

func connectToDB(app *Config) *sql.DB {
	// connect to the database
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			app.ErrorLog.Println("Could not connect to the database: %s", err)
		} else {
			app.InfoLog.Println("Connected to the database successfully!")
			return connection
		}

		if counts > 10 {
			return nil
		}

		app.InfoLog.Println("Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
		counts++
		continue

	}

}

func openDB(dsn string) (*sql.DB, error) {
	// open the database connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {

	gob.Register(data.User{})
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}

	return redisPool
}

func (app *Config) listenForShutdown() {

	// listen for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.shutdown()
	os.Exit(0)

}

func (app *Config) shutdown() {
	app.InfoLog.Println("Shutting down the server")

	// block until cleanup is complete
	app.Wait.Wait()

	// close the channel
	app.Mailer.DoneChan <- true
	app.ErrorChanDone <- true

	app.InfoLog.Println("Server shutdown complete, closing channels...")

	close(app.Mailer.ErrorChan)
	close(app.Mailer.DoneChan)
	close(app.Mailer.MailerChan)
	close(app.ErrorChan)
	close(app.ErrorChanDone)

}

func (app *Config) CreateMail() Mail {

	errorChan := make(chan error)
	mailerChan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:      os.Getenv("DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        mailPort,
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailerChan,
		DoneChan:    mailerDoneChan,
	}

	return m
}
