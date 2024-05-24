package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/joho/godotenv"
)

var ctx = context.Background()
var err error

func main() {
	err = godotenv.Load()
	PORT := os.Getenv("PORT")
	check(err, "error in loading env")

	http.HandleFunc("/ador/prod/set", setPasswordOracleAdor)
	http.HandleFunc("/ador/prod/get", getPasswordOracleAdor)

	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
	check(err, "error in serving application")
	server := fmt.Sprintf("listining to port %s", PORT)
	log.Println(server)
}

func db() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	// pong, err := rdb.Ping(ctx).Result()
	// check(err, "error in getting value")
	return rdb
}

func check(err error, s string) {
	if err != nil {
		log.Println(s, err)
	}
}

func setPasswordOracleAdor(res http.ResponseWriter, req *http.Request) {
	api_key := req.Header.Get("api_key")
	if api_key == os.Getenv("API_KEY") {

		redisIntance := db()
		password := req.URL.Query().Get("password")
		log.Println(password)

		err = redisIntance.Set(ctx, "oracle-ador-password", password, 0).Err()
		check(err, "error in setting value")
		log.Println("Value Updated")
	} else {
		json.NewEncoder(res).Encode("error")
	}

}

func getPasswordOracleAdor(res http.ResponseWriter, req *http.Request) {
	api_key := req.Header.Get("api_key")
	if api_key == os.Getenv("API_KEY") {
		redisIntance := db()
		val, err := redisIntance.Get(ctx, "oracle-ador-password").Result()
		check(err, "error in getting redis value")
		log.Println(val)
		json.NewEncoder(res).Encode(val)
	} else {
		json.NewEncoder(res).Encode("error")
	}
}
