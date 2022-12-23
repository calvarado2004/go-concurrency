package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/calvarado2004/go-concurrency/data"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initDB()

	// create web sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create wait group
	wg := sync.WaitGroup{}

	// set up application config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Models:   data.New(db),
	}

	// set up email

	// listen for signals
	go app.listenForShutdown()

	// listen for web connections
	app.serve()

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

func initDB() *sql.DB {
	// connect to the database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to the database")
	}

	return conn
}

func connectToDB() *sql.DB {
	// connect to the database
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("Could not connect to the database: %s", err)
		} else {
			log.Printf("Connected to the database")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Printf("Backing off and trying again")
		time.Sleep(1 * time.Second)
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

	app.Shutdown()
	os.Exit(0)

}

func (app *Config) Shutdown() {
	app.InfoLog.Println("Shutting down the server")

	// block until cleanup is complete
	app.Wait.Wait()

	app.InfoLog.Println("Server shutdown complete, closing channels...")

}
