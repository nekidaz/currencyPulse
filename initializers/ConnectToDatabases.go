package initializers

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func ConnectToPostgres() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect PostgresSQL :", err)
	}
}

func ConnectToRedis() {
	dsn := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	RDB = redis.NewClient(&redis.Options{
		Addr: dsn, // Поменяйте адрес и порт на соответствующие вашей конфигурации Redis.
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
}
