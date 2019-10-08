package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func readLines(fpath string) []string {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatalf("Can't find %s", fpath)
	}
	lines := strings.Split(strings.Trim(string(content), "\n"), "\n")
	return lines
}

func main() {
	loaded := godotenv.Load("../../.env")
	if loaded != nil {
		log.Println("No .env file found; using defaults")
	}

	prefix, found := os.LookupEnv("REDIS_PREFIX")
	if !found {
		prefix = "PB"
	}
	rkey := fmt.Sprintf("%s:YELLS", prefix)

	address, found := os.LookupEnv("REDIS_ADDRESS")
	if !found {
		address = "127.0.0.1:6379"
	}
	log.Printf("using redis @ %s to store our data", address)

	db := redis.NewClient(&redis.Options{Addr: address})

	seeds := readLines("SEEDS")

	pipe := db.Pipeline()
	for _, seed := range seeds {
		pipe.SAdd(rkey, seed)
	}
	_, err := pipe.Exec()
	if err != nil {
		log.Println("Could not write to redis!")
		log.Fatal(err)
	}

	log.Printf("Added %d shouts to the database at %s\n", len(seeds), rkey)
}
