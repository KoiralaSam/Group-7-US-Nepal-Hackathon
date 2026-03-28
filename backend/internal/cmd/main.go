package main

import (
	"log"
	"os"

	"github.com/KoiralaSam/Mindcare/backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	for _, p := range []string{"../.env", "../../.env", ".env"} {
		if err := godotenv.Load(p); err == nil {
			break
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set (add it to .env or export it)")
	}

	sqlDB, err := db.NewDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	log.Println("database connection OK")
}
