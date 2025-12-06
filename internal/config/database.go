package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() *pgx.Conn {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Одно из значений не задано")
	}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}

	log.Println("Подключение к базе прошло успешно!")
	return conn
}
