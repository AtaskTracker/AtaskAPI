package main

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/server"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := godotenv.Load("secrets/.env")
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		os.Getenv("MONGO_URI"),
	))
	if err != nil {
		log.Fatal(err)
	}
	redis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if res := redis.Ping(context.Background()); res.Err() != nil {
		log.Fatal("redis failed: ", res.Err())
	}
	server := server.NewServer(client, redis)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	server.Start(":" + port)
}
