package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	db       = ""
	duration = time.Minute * 15
	sleep    = time.Second * 1
)

func init() {
	flag.StringVar(&db, "db", db, "PostgreSQL database URI to connect (required)")
	flag.DurationVar(&duration, "duration", duration, "how long we should wait before returning an error")
	flag.DurationVar(&sleep, "sleep", sleep, "how long we should wait before next try")
}

func main() {
	flag.Parse()
	if db == "" {
		fmt.Printf("Error: db - is required option.\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	log.Println("start to wait for the database...")
	for {
		select {
		case <-ctx.Done():
			log.Fatalf("waiting failed: %s", ctx.Err())
		case <-time.After(sleep):
		}
		db, err := sql.Open("pgx", db)
		if err != nil {
			log.Println("can't connect to database...")
			continue
		}
		err = db.Ping()
		if err != nil {
			log.Println("can't ping database...")
			continue
		}
		log.Println("connected!")
		break
	}
}
