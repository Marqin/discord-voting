package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db   *bolt.DB
	conf *Config
)

func main() {
	var err error
	db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	conf, err = getConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = startBot(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
}
